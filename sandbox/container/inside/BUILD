cc_binary(
    name = "omogenjudge-sandboxr",
    srcs = ["init.cc"],
    visibility = [
        "//sandbox/container:__subpackages__",
        "//sandbox/deb:__subpackages__",
    ],
    deps = [
        ":inside",
        "@com_github_google_glog//:glog",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "inside",
    srcs = [
        "quota.cc",
        "chroot.cc",
        "setup.cc",
    ],
    hdrs = [
        "quota.h",
        "chroot.h",
        "setup.h",
    ],
    visibility = [
        "//sandbox/container/outside:__subpackages__",
    ],
    deps = [
        "@com_github_google_glog//:glog",
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
