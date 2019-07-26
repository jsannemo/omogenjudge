load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "omogenrunner_cli",
    embed = [":go_default_library"],
    visibility = ["//runner/deb:__subpackages__"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jsannemo/omogenjudge/runner/cli",
    visibility = ["//visibility:private"],
    deps = [
        "//runner/api:go_default_library",
        "//runner/client:go_default_library",
        "@com_github_google_logger//:go_default_library",
    ],
)
