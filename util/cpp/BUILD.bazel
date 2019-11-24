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
        "@com_github_gflags_gflags//:gflags",
        "@com_github_google_glog//:glog",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "test_main",
    srcs = ["test.cc"],
    linkopts = ["-pthread"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_google_glog//:glog",
        "@com_github_google_googletest//:gtest",
    ],
)
