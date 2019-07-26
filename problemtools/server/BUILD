load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "omogentools_server",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/server",
    visibility = ["//visibility:private"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/parser:go_default_library",
        "//problemtools/verifier:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
