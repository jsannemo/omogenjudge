load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "problems",
    srcs = [
        "list.go",
        "submit.go",
        "view.go",
    ],
    data = [":templates"],
    importpath = "github.com/jsannemo/omogenjudge/frontend/problems",
    visibility = ["//frontend/server:__subpackages__"],
    deps = [
        "//frontend/paths",
        "//frontend/request",
        "//storage/models:go_default_library",
        "//storage/problems",
        "//storage/submissions:go_default_library",
        "@com_github_gorilla_mux//:go_default_library",
    ],
)

filegroup(
    name = "templates",
    srcs = glob(["*.tpl"]),
)
