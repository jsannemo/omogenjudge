cc_binary(
    name = "omogenexec_cli",
    srcs = ["main.cc"],
    visibility = ["//exec/deb:__subpackages__"],
    deps = [
        "//exec/api:exec_cc_proto",
        "//external:gflags",
        "//external:glog",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
