package main

//go:generate go run ./internal/build-tools/embed.go -source=./configs/platform_cfg.yaml -target=./internal/platform/platform_cfg.go -name=PlatformConfig -package=platform
//go:generate go run ./internal/build-tools/embed.go -source=./configs/commands_cfg.yaml -target=./internal/builders/commands_cfg.go -name=CommandsConfig -package=builders
//go:generate go run ./internal/build-tools/embed.go -source=./configs/version.yaml -target=./internal/version/version.go -name=VersionConfig -package=version

