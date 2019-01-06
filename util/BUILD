cc_library(
    name = "util",
    srcs = [
        "files.cc",
        "time.cc",
    ],
    hdrs = [
        "files.h",
        "statusor.h",
        "time.h",
    ],
    visibility = ["//visibility:public"],
    deps = [
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
