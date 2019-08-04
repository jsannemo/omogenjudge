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
    tag = "v3.8.0",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

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
    sha256 = "8df59f11fb697743cbb3f26cfb8750395f30471e9eabde0d174c3aebc7a1cd39",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.19.1/rules_go-0.19.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.19.1/rules_go-0.19.1.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

# Gazelle, BUILD file generator for go

http_archive(
    name = "bazel_gazelle",
    sha256 = "be9296bfd64882e3c08e3283c58fcb461fa6dd3c171764fcc4cf322f60615a9b",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/0.18.1/bazel-gazelle-0.18.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.18.1/bazel-gazelle-0.18.1.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# Proto/gRPC

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    strip_prefix = "rules_proto-d9a123032f8436dbc34069cfc3207f2810a494ee",
    urls = ["https://github.com/stackb/rules_proto/archive/d9a123032f8436dbc34069cfc3207f2810a494ee.tar.gz"],
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
    commit = "1d89a3c832915b2314551c1d2a506874d62e53f7",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "org_golang_x_net",
    commit = "65e2d4e15006aab9813ff8769e768bbf4bb667a0",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_text",
    commit = "e6919f6577db79269a6443b9dc46d18f2238fb5d",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "4def268fd1a49955bfb3dda92fe3db4f924f2285",
    importpath = "golang.org/x/crypto",
)

# Markdown for go
go_repository(
    name = "com_github_gomarkdown",
    commit = "ee6a7931a1e4b802c9ff93e4dabcabacf4cb91db",
    importpath = "github.com/gomarkdown/markdown",
)

# Yaml for go
go_repository(
    name = "in_gopkg_yaml",
    commit = "51d6538a90f86fe93ac480b35f37b2be17fef232",
    importpath = "gopkg.in/yaml.v2",
)

# Postgres for go
go_repository(
    name = "com_github_lib_pq",
    commit = "3427c32cb71afc948325f299f040e53c1dd78979",
    importpath = "github.com/lib/pq",
)

go_repository(
    name = "com_github_jmoiron_sqlx",
    commit = "38398a30ed8516ffda617a04c822de09df8a3ec5",
    importpath = "github.com/jmoiron/sqlx",
)

# Go template helpers
go_repository (
    name = "com_github_masterminds_sprig",
    commit = "7525b3376b8792ab24d07381324e4e2463e3356b",
    importpath = "github.com/Masterminds/sprig",
)

go_repository (
    name = "com_github_masterminds_goutils",
    commit = "41ac8693c5c10a92ea1ff5ac3a7f95646f6123b0",
    importpath = "github.com/Masterminds/goutils",
)

go_repository (
    name = "com_github_masterminds_semver",
    commit = "0fd41f6ff0825cf7efae00e706120bdd48914d93",
    importpath = "github.com/Masterminds/semver",
)

go_repository(
    name = "com_github_imdario_mergo",
    commit = "4c317f2286be3bd0c4f1a0e622edc6398ec4656d",
    importpath = "github.com/imdario/mergo",
)

go_repository(
    name = "com_github_google_uuid",
    commit = "c2e93f3ae59f2904160ceaab466009f965df46d6",
    importpath = "github.com/google/uuid",
)

go_repository(
    name = "com_github_huandu_xstrings",
    commit = "8bbcf2f9ccb55755e748b7644164cd4bdce94c1d",
    importpath = "github.com/huandu/xstrings",
)

# Gorilla utilities

go_repository(
    name = "com_github_gorilla_securecookie",
    importpath = "github.com/gorilla/securecookie",
    sum = "h1:miw7JPhV+b/lAHSXz4qd/nN9jRiAFV5FwjeKyCS8BvQ=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_gorilla_sessions",
    commit = "4355a998706e83fe1d71c31b07af94e34f68d74a",
    importpath = "github.com/gorilla/sessions",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:gnP5JzjVOuiZD07fKKToCAOjS0yOpj/qPETTXCCS6hw=",
    version = "v1.7.3",
)

# Go logger

go_repository(
    name = "com_github_google_logger",
    commit = "7047ffcb7339f3f59be32de74a92217cb17cb40c",
    importpath = "github.com/google/logger",
)

# Go syntax highlighter

go_repository(
    name = "com_github_alecthomas_chroma",
    commit = "f8432cf78f68e5adf203ad5cefaaf6244650b4d1",
    importpath = "github.com/alecthomas/chroma",
)

go_repository(
    name = "com_github_danwakefield_fnmatch",
    importpath = "github.com/danwakefield/fnmatch",
    sum = "h1:y5HC9v93H5EPKqaS1UYVg1uYah5Xf51mBfIoWehClUQ=",
    version = "v0.0.0-20160403171240-cbb64ac3d964",
)

go_repository(
    name = "com_github_dlclark_regexp2",
    importpath = "github.com/dlclark/regexp2",
    sum = "h1:CqB4MjHw0MFCDj+PHHjiESmHX+N7t0tJzKvC6M97BRg=",
    version = "v1.1.6",
)

# Buildifier

http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-master",
    url = "https://github.com/bazelbuild/buildtools/archive/master.zip",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
