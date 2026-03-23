BINARY := assets
PLUGIN_DIR := $(HOME)/.gdt/plugins/gdt-plugin-assets
GO := go

.PHONY: build install clean test

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./...

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY) $(PLUGIN_DIR)/
	cp plugin.toml $(PLUGIN_DIR)/

clean:
	rm -f $(BINARY)
