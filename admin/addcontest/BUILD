load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["addcontest.go"],
    importpath = "github.com/jsannemo/omogenjudge/admin/addcontest",
    visibility = ["//visibility:private"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/client:go_default_library",
        "@com_github_google_logger//:go_default_library",
    ],
)

go_binary(
    name = "addcontest",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
