load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "github.com/jsannemo/omogenjudge/filehandler/client",
    visibility = ["//visibility:public"],
    deps = [
        "//filehandler/api:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
