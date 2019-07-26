cc_library(
    name = "outside",
    srcs = [
        "container.cc",
        "container_id.cc",
    ],
    hdrs = [
        "container.h",
        "container_id.h",
    ],
    visibility = ["//visibility:public"],
    deps = [
        ":internals",
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//sandbox/container/inside",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
    ],
)

cc_library(
    name = "internals",
    srcs = ["cgroups.cc"],
    hdrs = ["cgroups.h"],
    deps = [
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
