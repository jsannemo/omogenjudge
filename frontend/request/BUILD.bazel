load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "request",
    srcs = ["request.go"],
    importpath = "github.com/jsannemo/omogenjudge/frontend/request",
    visibility = ["//frontend:__subpackages__"],
    deps = [
        "//frontend/templates:go_default_library",
        "//storage/models:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@com_github_gorilla_sessions//:go_default_library",
        "@org_golang_x_text//language:go_default_library",
    ],
)
