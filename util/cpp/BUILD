cc_library(
    name = "cpp",
    srcs = [
        "config.cc",
        "files.cc",
        "time.cc",
    ],
    hdrs = [
        "config.h",
        "files.h",
        "statusor.h",
        "time.h",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//external:gflags",
        "//external:glog",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "test_main",
    srcs = ["test.cc"],
    linkopts = ["-pthread"],
    visibility = ["//visibility:public"],
    deps = [
        "//external:glog",
        "@gtest",
    ],
)
