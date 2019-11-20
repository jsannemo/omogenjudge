load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@rules_proto_grpc//go:go_grpc_library.bzl", "go_grpc_library")

package(default_visibility = ["//visibility:public"])

go_proto_library(
    name = "omogen_problemtools_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/api",
    proto = ":omogen_problemtools_proto",
    deps = ["//runner/api:go_default_library"],
)

go_library(
    name = "go_default_library",
    embed = [":omogen_problemtools_go_proto"],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/api",
)

proto_library(
    name = "omogen_problemtools_proto",
    srcs = [
        "course.proto",
        "problem.proto",
        "tools.proto",
    ],
    deps = ["//runner/api:omogen_runner_proto"],
)
