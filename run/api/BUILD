package(default_visibility = ["//visibility:public"])

load("@build_stack_rules_proto//go:go_proto_library.bzl", "go_proto_library")
load("@build_stack_rules_proto//go:go_grpc_library.bzl", "go_grpc_library")

proto_library(
    name = "run_proto",
    srcs = ["run.proto"],
    deps =
        [
            ":language_proto",
            ":program_proto",
        ],
)

go_grpc_library(
    name = "run_go_proto",
    deps = [":run_proto"],
    importpath = "github.com/jsannemo/omogenexec/run/api",
)

proto_library(
    name = "language_proto",
    srcs = ["language.proto"],
)

go_proto_library(
    name = "language_go_proto",
    deps = [":language_proto"],
    importpath = "github.com/jsannemo/omogenexec/run/api",
)

proto_library(
    name = "program_proto",
    srcs = ["program.proto"],
    deps = [":language_proto"],
)

go_proto_library(
    name = "program_go_proto",
    deps = [":language_proto"],
    importpath = "github.com/jsannemo/omogenexec/run/api",
)
