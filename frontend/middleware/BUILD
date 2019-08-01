load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "middleware",
    srcs = [
        "i18n.go",
        "middlewares.go",
        "session.go",
        "user.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/frontend/middleware",
    visibility = ["//frontend:__subpackages__"],
    deps = [
        "//frontend/request",
        "//storage/users",
        "@com_github_gorilla_sessions//:go_default_library",
        "@org_golang_x_text//language:go_default_library",
    ],
)
