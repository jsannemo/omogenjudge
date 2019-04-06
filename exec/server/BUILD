cc_binary(
    name = "omogenexec_server",
    srcs = [
        "executor.cc",
        "executor.h",
        "main.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//exec/api:exec_cc_proto",
        "//exec/container/outside",
        "@com_google_absl//absl/base",
        "@com_google_absl//absl/synchronization",
    ],
)
