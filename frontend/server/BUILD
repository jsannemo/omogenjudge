load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "omogenfrontend_server",
    data = [
        "//frontend/resources",
        "//frontend/templates",
    ],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "processor.go",
        "router.go",
        "server.go",
        "static.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/frontend/server",
    visibility = ["//visibility:private"],
    deps = [
        "//frontend/middleware",
        "//frontend/paths",
        "//frontend/problems",
        "//frontend/request",
        "//frontend/submissions:go_default_library",
        "//frontend/users",
        "@com_github_google_logger//:go_default_library",
        "@com_github_gorilla_mux//:go_default_library",
    ],
)
