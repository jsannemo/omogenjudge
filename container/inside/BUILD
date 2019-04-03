cc_binary(
    name = "omogenexec_runner",
    srcs = ["init.cc"],
    visibility = [
        "//container:__subpackages__",
        "//deb:__subpackages__",
    ],
    deps = [
        ":inside",
        "//external:glog",
        "//proto:container_cc_proto",
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
    visibility = ["//container:__subpackages__"],
    deps = [
        "//api:omogenexec_cc_proto",
        "//external:glog",
        "//proto:container_cc_proto",
        "//util",
        "@com_google_absl//absl/strings",
    ],
)
