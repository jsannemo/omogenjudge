load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")
load("@build_stack_rules_proto//cpp:cpp_grpc_library.bzl", "cpp_grpc_library")
load("@build_stack_rules_proto//go:go_grpc_library.bzl", "go_grpc_library")

cpp_grpc_library(
    name = "omogen_sandbox_cpp_proto",
    visibility = ["//visibility:public"],
    deps = [":omogen_sandbox_proto"],
)

go_proto_library(
    name = "omogen_sandbox_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jsannemo/omogenjudge/sandbox/api",
    proto = ":omogen_sandbox_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    embed = [":omogen_sandbox_go_proto"],
    importpath = "github.com/jsannemo/omogenjudge/sandbox/api",
    visibility = ["//visibility:public"],
)

proto_library(
    name = "omogen_sandbox_proto",
    srcs = [
        "containerspec.proto",
        "exec.proto",
        "execspec.proto",
        "resources.proto",
    ],
    visibility = ["//visibility:public"],
)
