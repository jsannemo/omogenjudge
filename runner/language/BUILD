load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "cpp.go",
        "detect.go",
        "doc.go",
        "language.go",
        "python.go",
    ],
    importpath = "github.com/jsannemo/omogenjudge/runner/language",
    visibility = ["//visibility:public"],
    deps = [
        "//runner/api:go_default_library",
        "//runner/compilers:go_default_library",
        "//runner/runners:go_default_library",
        "//sandbox/api:go_default_library",
        "//util/go/commands:go_default_library",
        "//util/go/files:go_default_library",
        "//util/go/users:go_default_library",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
