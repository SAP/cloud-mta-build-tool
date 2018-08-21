package main

//go:generate go run ./cmd/build-tools/embed.go -source=./platform_cfg.yaml -target=./cmd/platform/platform_cfg.go -name=PlatformConfig -package=platform
//go:generate go run ./cmd/build-tools/embed.go -source=./commands_cfg.yaml -target=./cmd/builders/commands_cfg.go -name=CommandsConfig -package=builders
