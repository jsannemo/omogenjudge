cc_library(
    name = "util",
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
        "@com_google_absl//absl/strings",
        "//external:gflags",
        "//external:glog",
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
