load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "view.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/frontend/editor",
    visibility = ["//visibility:public"],
    deps = [
        "//frontend/request",
        "//frontend/util:go_default_library",
    ],
)
