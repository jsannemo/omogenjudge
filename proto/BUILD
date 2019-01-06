load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")

proto_library(
    name = "container_proto",
    srcs = ["container.proto"],
    deps = [
        "//api:execspec_proto",
        "//api:resources_proto",
    ],
)

cpp_proto_library(
    name = "container_cc_proto",
    visibility = [
        "//container:__subpackages__",
    ],
    deps = [":container_proto"],
)
