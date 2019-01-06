package(default_visibility = ["//visibility:public"])

load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")
load("@build_stack_rules_proto//cpp:cpp_grpc_library.bzl", "cpp_grpc_library")

proto_library(
    name = "omogenexec_proto",
    srcs = ["omogenexec.proto"],
    deps = [
        ":containerspec_proto",
        ":execspec_proto",
    ],
)

cpp_grpc_library(
    name = "omogenexec_cc_proto",
    deps = [":omogenexec_proto"],
)

proto_library(
    name = "execspec_proto",
    srcs = ["execspec.proto"],
    deps = [":resources_proto"],
)

cpp_proto_library(
    name = "execspec_cc_proto",
    deps = [":execspec_proto"],
)

proto_library(
    name = "containerspec_proto",
    srcs = ["containerspec.proto"],
)

proto_library(
    name = "resources_proto",
    srcs = ["resources.proto"],
)

cpp_proto_library(
    name = "resources_cc_proto",
    deps = [":resources_proto"],
)
