load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "github.com/jsannemo/omogenjudge/runner/client",
    visibility = ["//visibility:public"],
    deps = [
        "//runner/api:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
