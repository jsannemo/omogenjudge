package(default_visibility = ["//visibility:public"])

cc_proto_library(
        name = "omogenexec_cc_proto",
        deps = ["omogenexec_proto"],
)

proto_library(
        name = "omogenexec_proto",
        srcs = ["omogenexec.proto"],
        deps = [":execspec_proto"]
)

cc_proto_library(
        name = "execspec_cc_proto",
        deps = [":execspec_proto"]
)

proto_library(
        name = "execspec_proto",
        srcs = ["execspec.proto"],
        deps = [":resources_proto"]
)

cc_proto_library(
        name = "resources_cc_proto",
        deps = [":resources_proto"]
)

proto_library(
        name = "resources_proto",
        srcs = ["resources.proto"]
)
