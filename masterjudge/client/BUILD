load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "github.com/jsannemo/omogenjudge/masterjudge/client",
    visibility = ["//visibility:public"],
    deps = [
        "//masterjudge/api:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
