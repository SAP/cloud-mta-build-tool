package main

//go:generate go run ./internal/buildtools/embed.go -source=./configs/platform_cfg.yaml -target=./internal/platform/platform_cfg.go -name=PlatformConfig -package=platform
//go:generate go run ./internal/buildtools/embed.go -source=./configs/module_type_cfg.yaml -target=./internal/commands/module_type_cfg.go -name=ModuleTypeConfig -package=commands
//go:generate go run ./internal/buildtools/embed.go -source=./configs/content_type_cfg.yaml -target=./internal/conttype/content_type_cfg.go -name=ContentTypeConfig -package=conttype
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/manifest.txt -target=./internal/tpl/manifest.go -name=Manifest -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./configs/builder_type_cfg.yaml -target=./internal/commands/builder_type_cfg.go -name=BuilderTypeConfig -package=commands
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_args.txt -target=./internal/tpl/base_args.go -name=baseArgs -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_post.txt -target=./internal/tpl/base_post.go -name=basePost -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_post_dep.txt -target=./internal/tpl/base_post_dep.go -name=basePostDep -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_pre_default.txt -target=./internal/tpl/base_pre_default.go -name=basePreDefault -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_pre_default_dep.txt -target=./internal/tpl/base_pre_default_dep.go -name=basePreDefaultDep -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_pre_verbose.txt -target=./internal/tpl/base_pre_verbose.go -name=basePreVerbose -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/base_pre_verbose_dep.txt -target=./internal/tpl/base_pre_verbose_dep.go -name=basePreVerboseDep -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/make_default.txt -target=./internal/tpl/make_default.go -name=makeDefault -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/make_deployment.txt -target=./internal/tpl/make_deployment.go -name=makeDeployment -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/make_verbose.txt -target=./internal/tpl/make_verbose.go -name=makeVerbose -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./internal/tpl/make_verbose_dep.txt -target=./internal/tpl/make_verbose_dep.go -name=makeVerboseDep -package=tpl
//go:generate go run ./internal/buildtools/embed.go -source=./configs/version.yaml -target=./internal/version/version_cfg.go -name=VersionConfig -package=version
