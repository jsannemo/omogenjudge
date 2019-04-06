cc_binary(
    name = "omogenexec_runner",
    srcs = ["init.cc"],
    visibility = [
        "//exec/container:__subpackages__",
        "//exec/deb:__subpackages__",
    ],
    deps = [
        ":inside",
        "//exec/proto:container_cc_proto",
        "//external:glog",
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
        "//exec/container/outside:__subpackages__",
    ],
    deps = [
        "//exec/api:exec_cc_proto",
        "//exec/proto:container_cc_proto",
        "//external:glog",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
