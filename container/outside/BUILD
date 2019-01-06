cc_library(
    name = "outside",
    srcs = ["container.cc"],
    hdrs = ["container.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":internals",
        "//api:omogenexec_cc_proto",
        "//container/inside",
        "//proto:container_cc_proto",
        "//util",
        "@com_google_absl//absl/strings",
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
