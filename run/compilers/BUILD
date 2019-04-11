load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "compilers",
    srcs = [
        "compilers.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//exec/api:exec_go_proto",
        "//run/api:run_go_proto",
        ],
    importpath = "github.com/jsannemo/omogenexec/run/compilers",
)
