.PHONY: all server plugins fawe waterdogpe fern clean

all: server fern plugins

server:
	GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod go build -o dragonfly .

fern:
	GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod go build -o fern ./cmd/fern/

plugins: fawe waterdogpe

fawe:
	GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod go build -buildmode=plugin -o plugins/fastasyncworldedit.pl ./plugin-src/fastasyncworldedit/

waterdogpe:
	GONOSUMCHECK=* GONOSUMDB=* GOFLAGS=-mod=mod go build -buildmode=plugin -o plugins/waterdogpe.pl ./plugin-src/waterdogpe/

clean:
	rm -f dragonfly plugins/*.pl
