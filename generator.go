package main

//go:generate go run ./internal/build-tools/embed.go -source=./platform_cfg.yaml -target=./internal/platform/platform_cfg.go -name=PlatformConfig -package=platform
//go:generate go run ./internal/build-tools/embed.go -source=./commands_cfg.yaml -target=./internal/builders/commands_cfg.go -name=CommandsConfig -package=builders
