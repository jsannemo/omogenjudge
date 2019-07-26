load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["addproblem.go"],
    importpath = "github.com/jsannemo/omogenjudge/admin/addproblem",
    visibility = ["//visibility:private"],
    deps = [
        "//problemtools/api:go_default_library",
        "//problemtools/client:go_default_library",
        "//storage/files:go_default_library",
        "//storage/problems",
        "//util/go/cli:go_default_library",
        "//util/go/filestore:go_default_library",
        "@com_github_google_logger//:go_default_library",
    ],
)

go_binary(
    name = "addproblem",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
