package server

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Origin-Net/FernMC/server/cmd"
)

type PluginMeta struct {
	Name         string
	Version      string
	Description  string
	Authors      []string
	Website      string
	Dependencies []string
}

type PluginContext struct {
	DataFolder string
	Logger     *slog.Logger
	Server     *Server
}

type Plugin interface {
	Meta() PluginMeta
	OnLoad(ctx PluginContext)
	OnEnable()
	OnDisable()
	OnUnload()
}

var (
	pluginMu   sync.Mutex
	plugins    = map[string]*loadedPlugin{}
	plugDir    = "plugins"
	plugSrcDir = "plugin-src"
	activeSrv  *Server
)

var pluginExts = []string{".fpl", ".so", ".pl"}

func pluginExt(name string) string {
	for _, ext := range pluginExts {
		if strings.HasSuffix(name, ext) {
			return ext
		}
	}
	return ""
}

func stripExt(name string) string {
	for _, ext := range pluginExts {
		if strings.HasSuffix(name, ext) {
			return name[:len(name)-len(ext)]
		}
	}
	return name
}

type loadedPlugin struct {
	Path     string
	Plugin   Plugin
	Meta     PluginMeta
	Context  PluginContext
	commands []string 
}

type pluginNode struct {
	plugin Plugin
	meta   PluginMeta
	path   string
}

func SetServer(srv *Server) {
	activeSrv = srv
}

func GetServer() *Server {
	return activeSrv
}

func SetPluginDir(dir string) {
	pluginMu.Lock()
	defer pluginMu.Unlock()
	plugDir = dir
}

func init() {
	_ = os.MkdirAll(plugDir, 0755)
	_ = os.MkdirAll(plugSrcDir, 0755)
}

func openPlugin(path string) (Plugin, PluginMeta, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, PluginMeta{}, fmt.Errorf("open %s: %w", path, err)
	}
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, PluginMeta{}, fmt.Errorf("lookup 'Plugin' in %s: %w", path, err)
	}

	if pp, ok := sym.(Plugin); ok {
		return pp, pp.Meta(), nil
	}

	v := reflect.ValueOf(sym)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, PluginMeta{}, fmt.Errorf("'Plugin' in %s is nil", path)
	}
	m := v.MethodByName("Meta")
	if !m.IsValid() {
		return nil, PluginMeta{}, fmt.Errorf("'Plugin' in %s is missing Meta() method (type: %s)", path, v.Type())
	}
	meta := m.Call(nil)[0].Interface().(PluginMeta)

	return &pluginReflect{val: v}, meta, nil
}

type pluginReflect struct {
	val reflect.Value
}

func (r *pluginReflect) Meta() PluginMeta {
	return r.val.MethodByName("Meta").Call(nil)[0].Interface().(PluginMeta)
}

func (r *pluginReflect) OnLoad(ctx PluginContext) {
	r.val.MethodByName("OnLoad").Call([]reflect.Value{reflect.ValueOf(ctx)})
}

func (r *pluginReflect) OnEnable() {
	r.val.MethodByName("OnEnable").Call(nil)
}

func (r *pluginReflect) OnDisable() {
	r.val.MethodByName("OnDisable").Call(nil)
}

func (r *pluginReflect) OnUnload() {
	r.val.MethodByName("OnUnload").Call(nil)
}

func resolveDeps(nodes []*pluginNode) ([]*pluginNode, error) {
	if len(nodes) == 0 {
		return nodes, nil
	}

	byName := map[string]*pluginNode{}
	for _, n := range nodes {
		name := strings.ToLower(n.meta.Name)
		if _, dup := byName[name]; dup {
			return nil, fmt.Errorf("duplicate plugin name: %s", name)
		}
		byName[name] = n
	}

	inDegree := map[string]int{}
	dependents := map[string][]string{}

	for _, n := range nodes {
		name := strings.ToLower(n.meta.Name)
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
		for _, dep := range n.meta.Dependencies {
			dep = strings.ToLower(dep)
			if _, ok := byName[dep]; !ok {
				return nil, fmt.Errorf("plugin %s depends on %s which is not found", name, dep)
			}
			dependents[dep] = append(dependents[dep], name)
			inDegree[name]++
		}
	}

	var queue []string
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	var order []string
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		order = append(order, name)
		for _, dep := range dependents[name] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(order) != len(nodes) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	result := make([]*pluginNode, len(order))
	for i, name := range order {
		result[i] = byName[name]
	}
	return result, nil
}

func pluginContext(srv *Server, name string) PluginContext {
	dir := filepath.Join(plugDir, name)
	_ = os.MkdirAll(dir, 0755)
	return PluginContext{
		DataFolder: dir,
		Logger:     srv.conf.Log.With("plugin", name),
		Server:     srv,
	}
}

func cmdNames() map[string]bool {
	cmds := cmd.Commands()
	names := make(map[string]bool, len(cmds))
	for n := range cmds {
		names[n] = true
	}
	return names
}

func sourceNewerThan(name, plPath string) bool {
	absPlugSrcDir, _ := filepath.Abs(plugSrcDir)
	srcDir := filepath.Join(absPlugSrcDir, name)
	plInfo, err := os.Stat(plPath)
	if err != nil {
		return true
	}
	var found bool
	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().After(plInfo.ModTime()) {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return false
	}
	return found
}

func serverNewerThan(plPath string) bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exeInfo, err := os.Stat(exe)
	if err != nil {
		return false
	}
	plInfo, err := os.Stat(plPath)
	if err != nil {
		return true
	}
	return exeInfo.ModTime().After(plInfo.ModTime())
}

var goBinPath string

var errDifferentVersion = fmt.Errorf("plugin version mismatch")

func isVersionMismatch(err error) bool {
	return err != nil && strings.Contains(err.Error(), "built with a different version")
}

func ensureGo() error {
	if goBinPath != "" {
		return nil
	}
	if p, err := exec.LookPath("go"); err == nil {
		goBinPath = p
		return nil
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable: %w", err)
	}
	localDir := filepath.Join(filepath.Dir(exe), "go")
	localGo := filepath.Join(localDir, "bin", "go")
	if _, err := os.Stat(localGo); err == nil {
		goBinPath = localGo
		return nil
	}
	slog.Info("Go not found, downloading...")
	arch := runtime.GOARCH
	url := fmt.Sprintf("https://go.dev/dl/go1.26.4.linux-%s.tar.gz", arch)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download Go: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download Go: HTTP %d", resp.StatusCode)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("decompress Go: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("extract Go: %w", err)
		}
		rel := strings.TrimPrefix(hdr.Name, "go/")
		if rel == "" {
			continue
		}
		target := filepath.Join(localDir, rel)
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(target), 0755)
		f, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("create %s: %w", target, err)
		}
		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return fmt.Errorf("write %s: %w", target, err)
		}
		f.Close()
		os.Chmod(target, os.FileMode(hdr.Mode))
	}
	goBinPath = filepath.Join(localDir, "bin", "go")
	slog.Info("Go installed", "path", goBinPath)
	return nil
}

func injectReplaces(pluginGoMod string, orig []byte) (restore func()) {
	restore = func() {}
	content := orig
	if content == nil {
		var err error
		content, err = os.ReadFile(pluginGoMod)
		if err != nil {
			return
		}
	}
	exe, exeErr := os.Executable()
	if exeErr != nil {
		return
	}
	srvDir := filepath.Dir(exe)
	srvMod, err := os.ReadFile(filepath.Join(srvDir, "go.mod"))
	if err != nil {
		return
	}
	var replaces []string
	for _, line := range strings.Split(string(srvMod), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "replace ") {
			continue
		}
		parts := strings.SplitN(trimmed, "=>", 2)
		if len(parts) != 2 {
			continue
		}
		target := strings.TrimSpace(parts[1])
		if strings.HasPrefix(target, "./") || strings.HasPrefix(target, "../") {
			target = filepath.Clean(filepath.Join(srvDir, target))
		}
		replaces = append(replaces, fmt.Sprintf("replace %s => %s", strings.TrimSpace(parts[0]), target))
	}
	if len(replaces) == 0 {
		return
	}
	modified := string(content)
	if !strings.HasSuffix(modified, "\n") {
		modified += "\n"
	}
	for _, r := range replaces {
		if !strings.Contains(modified, r) {
			modified += r + "\n"
		}
	}
	if modified == string(content) {
		return
	}
	restore = func() { os.WriteFile(pluginGoMod, content, 0644) }
	os.WriteFile(pluginGoMod, []byte(modified), 0644)
	return
}

func buildPlugin(srcDir, output string) error {
	pluginGoMod := filepath.Join(srcDir, "go.mod")

	origGoMod, _ := os.ReadFile(pluginGoMod)
	if origGoMod == nil {
		tmpMod := fmt.Sprintf("module __fern_plugin_%s\n\ngo 1.26\n\nrequire github.com/Origin-Net/FernMC v1.0.0\n", filepath.Base(srcDir))
		os.WriteFile(pluginGoMod, []byte(tmpMod), 0644)
		defer func() {
			os.Remove(pluginGoMod)
			os.Remove(filepath.Join(srcDir, "go.sum"))
		}()
		injectReplaces(pluginGoMod, []byte(tmpMod))
	} else {
		restore := injectReplaces(pluginGoMod, origGoMod)
		if restore != nil {
			defer restore()
		}
	}

	if err := ensureGo(); err != nil {
		return fmt.Errorf("ensure Go: %w", err)
	}

	tidy := exec.Command(goBinPath, "mod", "tidy")
	tidy.Dir = srcDir
	_ = tidy.Run()

	cmd := exec.Command(goBinPath, "build", "-buildmode=plugin", "-o", output, ".")
	cmd.Dir = srcDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go build failed: %w\n%s", err, string(out))
	}
	return nil
}

func safeOnLoad(p Plugin, ctx PluginContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	p.OnLoad(ctx)
	return nil
}

func safeOnEnable(p Plugin) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	p.OnEnable()
	return nil
}

func safeOnDisable(p Plugin) {
	defer func() { recover() }()
	p.OnDisable()
}

func safeOnUnload(p Plugin) {
	defer func() { recover() }()
	p.OnUnload()
}

func LoadPlugins(srv *Server) {
	_ = os.MkdirAll(plugDir, 0755)
	_ = os.MkdirAll(plugSrcDir, 0755)

	var nodes []*pluginNode

	absPlugDir, _ := filepath.Abs(plugDir)
	buildAndLoad := func(name, srcDir string) (Plugin, PluginMeta, error) {
		tmpPath := filepath.Join(absPlugDir, "."+name+".build.pl")
		if err := buildPlugin(srcDir, tmpPath); err != nil {
			return nil, PluginMeta{}, err
		}
		pl, meta, err := openPlugin(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			return nil, PluginMeta{}, err
		}
		finalPath := filepath.Join(absPlugDir, name+".pl")
		if err := os.Rename(tmpPath, finalPath); err != nil {
			os.Remove(tmpPath)
			return nil, PluginMeta{}, fmt.Errorf("rename: %w", err)
		}
		return pl, meta, nil
	}

	builtFromSource := map[string]bool{}
	srcEntries, _ := os.ReadDir(plugSrcDir)
	for _, s := range srcEntries {
		if !s.IsDir() {
			continue
		}
		name := s.Name()
		plPath := filepath.Join(absPlugDir, name+".pl")

		needsBuild := false
		if _, err := os.Stat(plPath); err != nil {
			needsBuild = true
		} else if sourceNewerThan(name, plPath) {
			needsBuild = true
		} else if serverNewerThan(plPath) {
			needsBuild = true
		}

		if !needsBuild {
			continue
		}

		srcDir := filepath.Join(plugSrcDir, name)
		srv.conf.Log.Info("Building plugin from source...", "name", name)
		pl, meta, err := buildAndLoad(name, srcDir)
		if err != nil {
			srv.conf.Log.Error(fmt.Sprintf("build/load failed for %s: %v", name, err))
			continue
		}
		srv.conf.Log.Info("Plugin built and loaded from source", "name", name, "version", meta.Version)
		builtFromSource[name] = true
		nodes = append(nodes, &pluginNode{plugin: pl, meta: meta, path: plPath})
	}

	entries, err := os.ReadDir(plugDir)
	if err != nil {
		srv.conf.Log.Error("read plugins dir: " + err.Error())
		return
	}

	for _, e := range entries {
		ext := pluginExt(e.Name())
		if e.IsDir() || ext == "" {
			continue
		}
		name := stripExt(e.Name())
		if builtFromSource[name] {
			continue
		}
		path := filepath.Join(absPlugDir, e.Name())
		p, meta, err := openPlugin(path)
		if err != nil {
			absSrcDir, _ := filepath.Abs(plugSrcDir)
			srcDir := filepath.Join(absSrcDir, name)
			if _, statErr := os.Stat(srcDir); statErr != nil {
				srv.conf.Log.Error(fmt.Sprintf("stale %s (no source to rebuild): %v", e.Name(), err))
				continue
			}
			srv.conf.Log.Warn(fmt.Sprintf("stale %s, rebuilding...", e.Name()))
			p, meta, err = buildAndLoad(name, srcDir)
			if err != nil {
				srv.conf.Log.Error(fmt.Sprintf("rebuild failed for %s: %v", name, err))
				continue
			}
		}
		nodes = append(nodes, &pluginNode{plugin: p, meta: meta, path: path})
	}

	if len(nodes) == 0 {
		return
	}

	ordered, err := resolveDeps(nodes)
	if err != nil {
		srv.conf.Log.Error("plugin dependency resolution: " + err.Error())
		return
	}

	type loadedEntry struct {
		node *pluginNode
		ctx  PluginContext
	}
	var loadedOrder []loadedEntry

	for _, n := range ordered {
		ctx := pluginContext(srv, n.meta.Name)
		if err := safeOnLoad(n.plugin, ctx); err != nil {
			srv.conf.Log.Error(fmt.Sprintf("plugin %s OnLoad panicked: %v", n.meta.Name, err))
			continue
		}
		pluginMu.Lock()
		plugins[n.meta.Name] = &loadedPlugin{
			Path:    n.path,
			Plugin:  n.plugin,
			Meta:    n.meta,
			Context: ctx,
		}
		pluginMu.Unlock()
		loadedOrder = append(loadedOrder, loadedEntry{node: n, ctx: ctx})
		srv.conf.Log.Info("Plugin loaded", "name", n.meta.Name, "version", n.meta.Version)
	}

	for _, l := range loadedOrder {
		before := cmdNames()
		if err := safeOnEnable(l.node.plugin); err != nil {
			srv.conf.Log.Error(fmt.Sprintf("plugin %s OnEnable panicked: %v", l.node.meta.Name, err))
			continue
		}
		after := cmdNames()
		var newCmds []string
		for n := range after {
			if !before[n] {
				newCmds = append(newCmds, n)
			}
		}
		pluginMu.Lock()
		if lp, ok := plugins[l.node.meta.Name]; ok {
			lp.commands = newCmds
		}
		pluginMu.Unlock()
		srv.conf.Log.Info("Plugin enabled", "name", l.node.meta.Name)
	}
}

func LoadPlugin(path string) (_ Plugin, loadErr error) {
	plug, meta, err := openPlugin(path)
	if err != nil {
		name := stripExt(filepath.Base(path))
		srcDir := filepath.Join(plugSrcDir, name)
		if _, statErr := os.Stat(srcDir); statErr != nil {
			return nil, fmt.Errorf("open failed: %w (no source in %s to rebuild)", err, srcDir)
		}
		absPlugDir, _ := filepath.Abs(plugDir)
		tmpPath := filepath.Join(absPlugDir, "."+name+".build.fpl")
		if buildErr := buildPlugin(srcDir, tmpPath); buildErr != nil {
			return nil, fmt.Errorf("open failed: %w; rebuild also failed: %v", err, buildErr)
		}
		plug, meta, err = openPlugin(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			return nil, fmt.Errorf("open after rebuild: %w", err)
		}
		finalPath := filepath.Join(absPlugDir, name+".fpl")
		if renameErr := os.Rename(tmpPath, finalPath); renameErr != nil {
			os.Remove(tmpPath)
			return nil, fmt.Errorf("rename rebuilt plugin: %w", renameErr)
		}
	}
	name := meta.Name

	pluginMu.Lock()
	if _, exists := plugins[name]; exists {
		pluginMu.Unlock()
		return nil, fmt.Errorf("plugin %s already loaded", name)
	}
	pluginMu.Unlock()

	srv := GetServer()
	if srv == nil {
		return nil, fmt.Errorf("server not available")
	}
	ctx := pluginContext(srv, name)

	if err := safeOnLoad(plug, ctx); err != nil {
		return nil, fmt.Errorf("plugin %s OnLoad failed: %w", name, err)
	}

	pluginMu.Lock()
	if _, exists := plugins[name]; exists {
		pluginMu.Unlock()
		safeOnUnload(plug)
		return nil, fmt.Errorf("plugin %s already loaded (race)", name)
	}
	plugins[name] = &loadedPlugin{
		Path:    path,
		Plugin:  plug,
		Meta:    meta,
		Context: ctx,
	}
	pluginMu.Unlock()

	before := cmdNames()
	if err := safeOnEnable(plug); err != nil {
		pluginMu.Lock()
		delete(plugins, name)
		pluginMu.Unlock()
		safeOnUnload(plug)
		return nil, fmt.Errorf("plugin %s OnEnable failed: %w", name, err)
	}
	after := cmdNames()
	var newCmds []string
	for n := range after {
		if !before[n] {
			newCmds = append(newCmds, n)
		}
	}
	pluginMu.Lock()
	if lp, ok := plugins[name]; ok {
		lp.commands = newCmds
	}
	pluginMu.Unlock()

	return plug, nil
}

func UnloadPlugin(name string) (Plugin, bool) {
	pluginMu.Lock()
	lp, ok := plugins[name]
	if !ok {
		pluginMu.Unlock()
		return nil, false
	}
	delete(plugins, name)
	pluginMu.Unlock()

	for _, c := range lp.commands {
		cmd.Unregister(c)
	}
	safeOnDisable(lp.Plugin)
	safeOnUnload(lp.Plugin)
	return lp.Plugin, true
}

func ReloadPlugin(name string, srv *Server) (_ Plugin, reloadErr error) {
	pluginMu.Lock()
	lp, ok := plugins[name]
	pluginMu.Unlock()
	if !ok {
		return nil, fmt.Errorf("plugin %s not loaded", name)
	}

	srcDir := filepath.Join(plugSrcDir, name)
	if _, statErr := os.Stat(srcDir); statErr != nil {
		return nil, fmt.Errorf("source directory for %s not found at %s", name, srcDir)
	}

	absPlugDir, _ := filepath.Abs(plugDir)
	tmpPath := filepath.Join(absPlugDir, "."+name+"-reload.fpl")
	if err := buildPlugin(srcDir, tmpPath); err != nil {
		return nil, fmt.Errorf("rebuild failed: %w", err)
	}
	defer os.Remove(tmpPath)

	newPlug, meta, err := openPlugin(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("open rebuilt plugin: %w", err)
	}

	ctx := pluginContext(srv, name)

	if err := safeOnLoad(newPlug, ctx); err != nil {
		return nil, fmt.Errorf("new %s OnLoad failed: %w", name, err)
	}

	for _, c := range lp.commands {
		cmd.Unregister(c)
	}
	safeOnDisable(lp.Plugin)
	safeOnUnload(lp.Plugin)

	pluginMu.Lock()
	plugins[name] = &loadedPlugin{
		Path:    lp.Path,
		Plugin:  newPlug,
		Meta:    meta,
		Context: ctx,
	}
	pluginMu.Unlock()

	before := cmdNames()
	if err := safeOnEnable(newPlug); err != nil {
		return nil, fmt.Errorf("new %s OnEnable failed: %w", name, err)
	}
	after := cmdNames()
	var newCmds []string
	for n := range after {
		if !before[n] {
			newCmds = append(newCmds, n)
		}
	}
	pluginMu.Lock()
	if lp, ok := plugins[name]; ok {
		lp.commands = newCmds
	}
	pluginMu.Unlock()

	finalPath := filepath.Join(absPlugDir, name+".fpl")
	if err := os.Rename(tmpPath, finalPath); err != nil {
		srv.conf.Log.Warn(fmt.Sprintf("could not update %s final path: %v", name, err))
	}

	return newPlug, nil
}

func ShutdownPlugins() {
	pluginMu.Lock()
	all := make([]*loadedPlugin, 0, len(plugins))
	for _, lp := range plugins {
		all = append(all, lp)
	}
	plugins = map[string]*loadedPlugin{}
	pluginMu.Unlock()

	for _, lp := range all {
		safeOnDisable(lp.Plugin)
		safeOnUnload(lp.Plugin)
	}
}

func PluginByName(name string) (Plugin, bool) {
	pluginMu.Lock()
	defer pluginMu.Unlock()
	lp, ok := plugins[name]
	if !ok {
		return nil, false
	}
	return lp.Plugin, true
}

func AllPlugins() []Plugin {
	pluginMu.Lock()
	defer pluginMu.Unlock()
	out := make([]Plugin, 0, len(plugins))
	for _, lp := range plugins {
		out = append(out, lp.Plugin)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Meta().Name < out[j].Meta().Name
	})
	return out
}
