cc_binary(
    name = "init",
    srcs = ["init.cc"],
    visibility = ["//container:__subpackages__"],
    deps = [
        ":inside",
        "//external:glog",
        "//proto:container_cc_proto",
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
