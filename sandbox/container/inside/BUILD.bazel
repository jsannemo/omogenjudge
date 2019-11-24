cc_binary(
    name = "omogenjudge-sandboxr",
    srcs = ["init.cc"],
    visibility = [
        "//sandbox/container:__subpackages__",
        "//sandbox/deb:__subpackages__",
    ],
    deps = [
        ":inside",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "@com_github_google_glog//:glog",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "inside",
    srcs = [
        "chroot.cc",
        "quota.cc",
        "setup.cc",
    ],
    hdrs = [
        "chroot.h",
        "quota.h",
        "setup.h",
    ],
    visibility = [
        "//sandbox/container/outside:__subpackages__",
    ],
    deps = [
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "//util/cpp",
        "@com_github_google_glog//:glog",
        "@com_google_absl//absl/strings",
    ],
)
