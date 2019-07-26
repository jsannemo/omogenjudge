load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "omogenrunner_server",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jsannemo/omogenjudge/runner/server",
    visibility = ["//visibility:private"],
    deps = [
        "//runner/api:go_default_library",
        "//runner/language:go_default_library",
        "//sandbox/api:go_default_library",
        "//sandbox/client:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
