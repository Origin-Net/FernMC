package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Origin-Net/FernMC/server"
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/pelletier/go-toml"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

func nowTag() string {
	return time.Now().Format("15:04:05")
}

func colorize(s string) string {
	return text.ANSI(s)
}

var levelColor = map[slog.Level]string{
	slog.LevelDebug: "\033[90m",
	slog.LevelInfo:  "\033[97m",
	slog.LevelWarn:  "\033[93m",
	slog.LevelError: "\033[91m",
}

var levelTag = map[slog.Level]string{
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARN",
	slog.LevelError: "ERROR",
}

type consoleHandler struct {
	w slog.Handler
}

func (h *consoleHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.w.Enabled(ctx, l)
}

func (h *consoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &consoleHandler{h.w.WithAttrs(attrs)}
}

func (h *consoleHandler) WithGroup(name string) slog.Handler {
	return &consoleHandler{h.w.WithGroup(name)}
}

func (h *consoleHandler) Handle(ctx context.Context, r slog.Record) error {
	tag, _ := levelTag[r.Level]
	c, _ := levelColor[r.Level]
	msg := colorize(r.Message)
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindString {
			attrs += fmt.Sprintf(" %s=\033[90m%s\033[0m", a.Key, a.Value.String())
		} else {
			attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		}
		return true
	})
	fmt.Printf("\033[97m[%s %s]\033[0m %s%s\033[0m%s\n", nowTag(), tag, c, msg, attrs)
	return nil
}

func printfInfo(format string, a ...any) {
	fmt.Printf("\033[97m[%s INFO]:\033[0m %s\033[0m\n", nowTag(), colorize(fmt.Sprintf(format, a...)))
}

func printfWarn(format string, a ...any) {
	fmt.Printf("\033[97m[%s WARN]:\033[0m \033[93m%s\033[0m\n", nowTag(), colorize(fmt.Sprintf(format, a...)))
}


type consoleSource struct{}

func (consoleSource) Name() string                            { return "Console" }
func (consoleSource) Position() mgl64.Vec3                     { return mgl64.Vec3{} }
func (consoleSource) SendCommandOutput(o *cmd.Output) {
	for _, m := range o.Messages() {
		fmt.Printf("\033[97m[%s INFO]:\033[0m %s\033[0m\n", nowTag(), colorize(m.String()))
	}
	for _, e := range o.Errors() {
		fmt.Printf("\033[97m[%s WARN]:\033[0m \033[91m%s\033[0m\n", nowTag(), colorize(e.Error()))
	}
}


type chatSubscriber struct{}

func (chatSubscriber) UUID() uuid.UUID { return uuid.New() }
func (chatSubscriber) Message(a ...any) {
	s := make([]string, len(a))
	for i, b := range a {
		s[i] = fmt.Sprint(b)
	}
	msg := colorize(text.ANSI(strings.Join(s, " ")))
	fmt.Printf("\033[97m[%s INFO]:\033[0m %s\033[0m\n", nowTag(), msg)
}

func main() {
	start := time.Now()
	fmt.Print("\033[H\033[2J")
	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.SetDefault(slog.New(&consoleHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})}))
	chat.Global.Subscribe(chatSubscriber{})

	envInfo := fmt.Sprintf("Running Go %s (%s; %s) on %s/%s (%s)",
		runtime.Version(), runtime.Compiler, runtime.Version(),
		runtime.GOOS, runtime.GOARCH, runtime.GOARCH,
	)
	printfInfo("[bootstrap] %s", envInfo)

	fernVer := "v1.01"
	fernCommit := "unknown"

	printfInfo("[bootstrap] Loading FernMC %s for Minecraft 1.26.30 (commit %s)", fernVer, fernCommit)

	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	go func() {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				fmt.Print("> ")
				continue
			}
			args := strings.SplitN(line, " ", 2)
			commandName := strings.ToLower(args[0])
			var argLine string
			if len(args) > 1 {
				argLine = args[1]
			}
			c, ok := cmd.ByAlias(commandName)
			if !ok {
				fmt.Printf("\033[97m[%s WARN]:\033[0m \033[91mUnknown command: %s\033[0m\n", nowTag(), commandName)
				fmt.Print("> ")
				continue
			}
			c.Execute(argLine, consoleSource{}, nil)
			fmt.Print("> ")
		}
	}()

	printfInfo("Done (%s)! For help, type \"help\"", time.Since(start).Round(time.Millisecond).String())

	for p := range srv.Accept() {
		slog.Info(fmt.Sprintf("%s[/127.0.0.1] logged in with entity id at %v", p.Name(), p.Position()))
	}
}



func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
