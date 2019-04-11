load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "omogenrun_server",
    srcs = [
        "main.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//exec/api:exec_go_proto",
        "//run/api:run_go_proto",
        "//run/language",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
