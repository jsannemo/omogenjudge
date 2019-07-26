load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "parser.go",
        "verify.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/testdata",
    visibility = ["//visibility:public"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/util:go_default_library",
        "@in_gopkg_yaml//:go_default_library",
    ],
)
