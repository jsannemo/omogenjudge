load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")

proto_library(
    name = "container_proto",
    srcs = ["container.proto"],
    deps = [
        "//exec/api:execspec_proto",
        "//exec/api:resources_proto",
    ],
)

cpp_proto_library(
    name = "container_cc_proto",
    visibility = [
        "//exec/container:__subpackages__",
    ],
    deps = [":container_proto"],
)
