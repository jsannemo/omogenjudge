load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "language",
    srcs = [
        "language.go",
        "python.go",
        ],
    visibility = ["//visibility:public"],
    deps = [
        "//exec/api:exec_go_proto",
        "//run/api:run_go_proto",
        "//run/compilers",
        "//run/runners",
        ],
    importpath = "github.com/jsannemo/omogenexec/run/language",
    )
