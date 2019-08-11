cc_binary(
    name = "omogenjudge-sandbox",
    srcs = [
        "executor.cc",
        "executor.h",
        "request_validator.h",
        "request_validator.cc",
        "main.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//sandbox/container/outside",
        "@com_google_absl//absl/base",
        "@com_google_absl//absl/synchronization",
    ],
)
