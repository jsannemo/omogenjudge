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
        "//api:omogenexec_cc_proto",
        "//container/inside",
        "//proto:container_cc_proto",
        "//util",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
    ],
)

cc_library(
    name = "internals",
    srcs = ["cgroups.cc"],
    hdrs = ["cgroups.h"],
    visibility = ["//container/outside:__pkg__"],
    deps = [
        "//api:omogenexec_cc_proto",
        "//util",
        "@com_google_absl//absl/strings",
    ],
)
