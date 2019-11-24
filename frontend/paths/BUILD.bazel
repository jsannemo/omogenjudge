load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "paths",
    srcs = ["paths.go"],
    importpath = "github.com/jsannemo/omogenjudge/frontend/paths",
    visibility = [
        "//frontend:__subpackages__",
        "//storage:__subpackages__",
    ],
    deps = ["@com_github_gorilla_mux//:go_default_library"],
)
