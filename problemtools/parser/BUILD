load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["parser.go"],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/parser",
    visibility = ["//visibility:public"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/metadata:go_default_library",
        "//problemtools/statement:go_default_library",
        "//problemtools/testdata:go_default_library",
        "//problemtools/util:go_default_library",
    ],
)
