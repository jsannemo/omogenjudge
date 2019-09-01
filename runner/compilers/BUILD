load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["compilers.go"],
    importpath = "github.com/jsannemo/omogenjudge/runner/compilers",
    visibility = ["//visibility:public"],
    deps = [
        "//runner/api:go_default_library",
        "//sandbox/api:go_default_library",
        "//util/go/files:go_default_library",
        "//util/go/users:go_default_library",
    ],
)
