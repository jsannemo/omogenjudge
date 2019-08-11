cc_binary(
    name = "omogensandbox_cli",
    srcs = ["main.cc"],
    visibility = ["//sandbox/deb:__subpackages__"],
    deps = [
        "@com_github_gflags_gflags//:gflags",
        "@com_github_google_glog//:glog",
        "//sandbox/api:omogen_sandbox_cpp_proto",
        "//util/cpp",
        "@com_google_absl//absl/strings",
    ],
)
