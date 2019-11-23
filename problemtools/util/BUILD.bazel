load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "reporter.go",
        "yaml.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/util",
    visibility = ["//visibility:public"],
    deps = ["@in_gopkg_yaml//:go_default_library"],
)
