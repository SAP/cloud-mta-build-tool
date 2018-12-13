package main

//go:generate go run ./internal/build-tools/embed.go -source=./configs/platform_cfg.yaml -target=./internal/platform/platform_cfg.go -name=PlatformConfig -package=platform
//go:generate go run ./internal/build-tools/embed.go -source=./configs/commands_cfg.yaml -target=./internal/commands/commands_cfg.go -name=CommandsConfig -package=commands
//go:generate go run ./internal/build-tools/embed.go -source=./configs/version.yaml -target=./internal/version/version_cfg.go -name=VersionConfig -package=version
//go:generate go run ./internal/build-tools/embed.go -source=./validations/schema.yaml -target=./validations/mta_schema.go -name=schemaDef -package=validate
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/base_post_default.txt -target=./internal/tpl/base_post_default.go -name=basePostDefault -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/base_post_verbose.txt -target=./internal/tpl/base_post_verbose.go -name=basePostVerbose -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/base_pre_default.txt -target=./internal/tpl/base_pre_default.go -name=basePreDefault -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/base_pre_verbose.txt -target=./internal/tpl/base_pre_verbose.go -name=basePreVerbose -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/make_default.txt -target=./internal/tpl/make_default.go -name=makeDefault -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/make_deployment.txt -target=./internal/tpl/make_deployment.go -name=makeDeployment -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/make_verbose.txt -target=./internal/tpl/make_verbose.go -name=makeVerbose -package=tpl
//go:generate go run ./internal/build-tools/embed.go -source=./internal/tpl/make_verbose_dep.txt -target=./internal/tpl/make_verbose_dep.go -name=makeVerboseDep -package=tpl
