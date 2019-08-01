load("@io_bazel_rules_go//go:def.bzl", "go_library")

filegroup(
    name = "templates",
    srcs = glob(["*.tpl"]),
    visibility = ["//frontend/server:__subpackages__"],
)

go_library(
    name = "go_default_library",
    srcs = ["templates.go"],
    importpath = "github.com/jsannemo/omogenjudge/frontend/templates",
    visibility = ["//visibility:public"],
)
