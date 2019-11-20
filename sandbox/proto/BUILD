load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto_grpc//cpp:cpp_proto_library.bzl", "cpp_proto_library")

cpp_proto_library(
    name = "omogen_sandbox_proto_cc_proto",
    visibility = ["//visibility:public"],
    deps = [":omogen_sandbox_proto_proto"],
)

proto_library(
    name = "omogen_sandbox_proto_proto",
    srcs = ["container.proto"],
    visibility = ["//visibility:public"],
    deps = ["//sandbox/api:omogen_sandbox_proto"],
)

go_proto_library(
    name = "omogen_sandbox_proto_go_proto",
    importpath = "github.com/jsannemo/omogenjudge/sandbox/proto",
    proto = ":omogen_sandbox_proto_proto",
    visibility = ["//visibility:public"],
    deps = ["//sandbox/api:go_default_library"],
)

go_library(
    name = "go_default_library",
    embed = [":omogen_sandbox_proto_go_proto"],
    importpath = "github.com/jsannemo/omogenjudge/sandbox/proto",
    visibility = ["//visibility:public"],
)
