cc_binary(
    name = "omogensandbox_cli",
    srcs = ["main.cc"],
    visibility = ["//sandbox/deb:__subpackages__"],
    deps = [
        "//external:gflags",
        "//external:glog",
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
