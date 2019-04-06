workspace(name = "omogenexec")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_github_gflags_gflags",
    remote = "https://github.com/gflags/gflags.git",
    tag = "v2.2.2",
)

bind(
    name = "gflags",
    actual = "@com_github_gflags_gflags//:gflags",
)

git_repository(
    name = "gtest",
    remote = "https://github.com/google/googletest.git",
    tag = "release-1.8.1",
)

git_repository(
    name = "com_google_protobuf",
    remote = "https://github.com/google/protobuf.git",
    tag = "v3.6.1.3",
)

git_repository(
    name = "com_google_absl",
    remote = "https://github.com/abseil/abseil-cpp.git",
    tag = "20180600",
)

git_repository(
    name = "com_github_glog_glog",
    commit = "0e4ce7c0c0f7cda7cc86017abd775cecf04074e0",
    remote = "https://github.com/google/glog.git",
)

bind(
    name = "glog",
    actual = "@com_github_glog_glog//:glog",
)

http_archive(
    name = "build_stack_rules_proto",
    strip_prefix = "rules_proto-f5d6eea6a4528bef3c1d3a44d486b51a214d61c2",
    urls = ["https://github.com/stackb/rules_proto/archive/f5d6eea6a4528bef3c1d3a44d486b51a214d61c2.tar.gz"],
)

load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_library")

cpp_proto_library()

load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_grpc_compile")

cpp_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f87fa87475ea107b3c69196f39c82b7bbf58fe27c62a338684c20ca17d1d8613",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.2/rules_go-0.16.2.tar.gz",
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-2a27d63db79086b75a7dd646cacce0e931535691",
    url = "https://github.com/bazelbuild/buildtools/archive/2a27d63db79086b75a7dd646cacce0e931535691.zip",
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

go_rules_dependencies()

go_register_toolchains()

buildifier_dependencies()
