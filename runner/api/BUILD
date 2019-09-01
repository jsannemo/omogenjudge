load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@com_github_stackb_rules_proto//go:go_grpc_library.bzl", "go_grpc_library")

package(default_visibility = ["//visibility:public"])

go_proto_library(
    name = "omogen_runner_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jsannemo/omogenjudge/runner/api",
    proto = ":omogen_runner_proto",
)

go_library(
    name = "go_default_library",
    embed = [":omogen_runner_go_proto"],
    importpath = "github.com/jsannemo/omogenjudge/runner/api",
)

proto_library(
    name = "omogen_runner_proto",
    srcs = [
        "eval.proto",
        "language.proto",
        "program.proto",
        "run_service.proto",
    ],
)
