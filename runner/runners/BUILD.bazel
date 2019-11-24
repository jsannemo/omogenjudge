load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "exec.go",
        "filelinker.go",
        "program.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/runner/runners",
    visibility = ["//visibility:public"],
    deps = [
        "//runner/api:go_default_library",
        "//sandbox/api:go_default_library",
        "//util/go/files:go_default_library",
        "//util/go/users:go_default_library",
        "@com_github_google_logger//:go_default_library",
    ],
)
