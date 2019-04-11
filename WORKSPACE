workspace(name = "omogenexec")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# Gflags

git_repository(
    name = "com_github_gflags_gflags",
    remote = "https://github.com/gflags/gflags.git",
    tag = "v2.2.2",
)

bind(
    name = "gflags",
    actual = "@com_github_gflags_gflags//:gflags",
)

# Gtest

git_repository(
    name = "gtest",
    remote = "https://github.com/google/googletest.git",
    tag = "release-1.8.1",
)

# Protobuf

git_repository(
    name = "com_google_protobuf",
    remote = "https://github.com/google/protobuf.git",
    tag = "v3.6.1.3",
)

# Abseil

git_repository(
    name = "com_google_absl",
    remote = "https://github.com/abseil/abseil-cpp.git",
    tag = "20180600",
)

# Glog

git_repository(
    name = "com_github_glog_glog",
    commit = "0e4ce7c0c0f7cda7cc86017abd775cecf04074e0",
    remote = "https://github.com/google/glog.git",
)

bind(
    name = "glog",
    actual = "@com_github_glog_glog//:glog",
)

# Golang

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.18.2/rules_go-0.18.2.tar.gz"],
    sha256 = "31f959ecf3687f6e0bb9d01e1e7a7153367ecd82816c9c0ae149cd0e5a92bf8c",
)
load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
)
load("@bazel_gazelle//:deps.bzl", "go_repository")


# Proto/gRPC

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/e783457abea020e7df6b94acb54f668c0473ae31.tar.gz"],
    strip_prefix = "rules_proto-e783457abea020e7df6b94acb54f668c0473ae31",
)

load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_library")
cpp_proto_library()
load("@build_stack_rules_proto//go:deps.bzl", "go_proto_library")
go_proto_library()

load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_grpc_library")
cpp_grpc_library()
load("@build_stack_rules_proto//go:deps.bzl", "go_grpc_library")
go_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
grpc_deps()

go_repository(
    name = "org_golang_google_grpc",
    commit = "3507fb8e1a5ad030303c106fef3a47c9fdad16ad",
    importpath = "google.golang.org/grpc",
)

# Buildifier

http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-2a27d63db79086b75a7dd646cacce0e931535691",
    url = "https://github.com/bazelbuild/buildtools/archive/2a27d63db79086b75a7dd646cacce0e931535691.zip",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
