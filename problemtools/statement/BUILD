load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "markdown.go",
        "parser.go",
        "verify.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/problemtools/statement",
    visibility = ["//visibility:public"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/util:go_default_library",
        "@com_github_gomarkdown//:go_default_library",
        "@com_github_gomarkdown//ast:go_default_library",
        "@com_github_gomarkdown//html:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_x_text//language:go_default_library",
    ],
)
