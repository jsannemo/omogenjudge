load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["runners.go"],
    importpath = "github.com/jsannemo/omogenjudge/runner/runners",
    visibility = ["//visibility:public"],
    deps = [
        "//runner/api:go_default_library",
        "//sandbox/api:go_default_library",
        "@com_github_google_logger//:go_default_library",
    ],
)
