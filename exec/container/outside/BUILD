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
        "//exec/api:exec_cc_proto",
        "//exec/container/inside",
        "//exec/proto:container_cc_proto",
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
        "//exec/api:exec_cc_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
