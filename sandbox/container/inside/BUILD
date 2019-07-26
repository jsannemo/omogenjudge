cc_binary(
    name = "omogensandbox_runner",
    srcs = ["init.cc"],
    visibility = [
        "//sandbox/container:__subpackages__",
        "//sandbox/deb:__subpackages__",
    ],
    deps = [
        ":inside",
        "//external:glog",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "inside",
    srcs = [
        "chroot.cc",
        "setup.cc",
    ],
    hdrs = [
        "chroot.h",
        "setup.h",
    ],
    visibility = [
        "//sandbox/container/outside:__subpackages__",
    ],
    deps = [
        "//external:glog",
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//sandbox/proto:omogen_sandbox_proto_cc_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
