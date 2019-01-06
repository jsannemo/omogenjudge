cc_binary(
    name = "daemon",
    srcs = [
        "executor.cc",
        "executor.h",
        "main.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//api:omogenexec_cc_proto",
        "//container/outside",
        "@com_google_absl//absl/base",
        "@com_google_absl//absl/synchronization",
    ],
)
