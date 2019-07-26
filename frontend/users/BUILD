load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "users",
    srcs = [
        "login.go",
        "logout.go",
        "register.go",
    ],
    data = [":templates"],
    importpath = "github.com/jsannemo/omogenjudge/frontend/users",
    visibility = ["//frontend/server:__subpackages__"],
    deps = [
        "//frontend/paths",
        "//frontend/request",
        "//storage/users",
    ],
)

filegroup(
    name = "templates",
    srcs = glob(["*.tpl"]),
)
