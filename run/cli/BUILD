load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "omogenrun_cli",
    srcs = ["main.go"],
    visibility = ["//run/deb:__subpackages__"],
    deps = [
        "//run/api:run_go_proto",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
