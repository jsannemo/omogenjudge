workspace(name = "omogenexec")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# ======== Build tools ========

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.20.2/rules_go-v0.20.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.20.2/rules_go-v0.20.2.tar.gz",
    ],
    sha256 = "b9aa86ec08a292b97ec4591cf578e020b35f98e12173bbd4a921f84f583aebd9",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

git_repository(
    name = "com_github_bazelbuild_buildtools",
    remote = "https://github.com/bazelbuild/buildtools.git",
    commit = "d7ccc5507c6c16e04f5e362e558d70b8b179b052",
    shallow_since = "1562930059 +0300",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")
buildifier_dependencies()

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

# ======== C++ Libraries ========

git_repository(
    name = "com_github_gflags_gflags",
    remote = "https://github.com/gflags/gflags.git",
    commit = "e171aa2d15ed9eb17054558e0b3a6a413bb01067",
    shallow_since = "1541971260 +0000",
)

git_repository(
    name = "com_github_google_googletest",
    remote = "https://github.com/google/googletest.git",
    commit = "2fe3bd994b3189899d93f1d5a881e725e046fdc2",
    shallow_since = "1535728917 -0400",
)

git_repository(
    name = "com_google_absl",
    remote = "https://github.com/abseil/abseil-cpp.git",
    commit = "aa844899c937bde5d2b24f276b59997e5b668bde",
    shallow_since = "1565288385 -0400"
)

git_repository(
    name = "com_github_google_glog",
    commit = "0e4ce7c0c0f7cda7cc86017abd775cecf04074e0",
    shallow_since = "1546302462 +0100",
    remote = "https://github.com/google/glog.git",
)

# ======== Protobuf ========
git_repository(
    name = "rules_proto_grpc",
    commit = "65b3876bf833fe72049217312727fd04e04e6a1e",
    remote = "https://github.com/rules-proto-grpc/rules_proto_grpc.git",
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos")
rules_proto_grpc_repos()

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains")
rules_proto_grpc_toolchains()

load("@rules_proto_grpc//cpp:repositories.bzl", rules_proto_grpc_cpp_repos="cpp_repos")
rules_proto_grpc_cpp_repos()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos="go_repos")
rules_proto_grpc_go_repos()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
grpc_deps()


# ======== Golang libraries ========

# gRPC
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

# Markdown parser
go_repository(
    name = "com_github_gomarkdown",
    commit = "ee6a7931a1e4b802c9ff93e4dabcabacf4cb91db",
    importpath = "github.com/gomarkdown/markdown",
)

# Yaml parser
go_repository(
    name = "in_gopkg_yaml",
    commit = "51d6538a90f86fe93ac480b35f37b2be17fef232",
    importpath = "gopkg.in/yaml.v2",
)

# Postgres driver
go_repository(
    name = "com_github_lib_pq",
    commit = "3427c32cb71afc948325f299f040e53c1dd78979",
    importpath = "github.com/lib/pq",
)

# More useful SQL library
go_repository(
    name = "com_github_jmoiron_sqlx",
    commit = "38398a30ed8516ffda617a04c822de09df8a3ec5",
    importpath = "github.com/jmoiron/sqlx",
)

# HTML template helpers
go_repository (
    name = "com_github_masterminds_sprig",
    commit = "7525b3376b8792ab24d07381324e4e2463e3356b",
    importpath = "github.com/Masterminds/sprig",
)

# Indirect dependency
go_repository (
    name = "com_github_masterminds_goutils",
    commit = "41ac8693c5c10a92ea1ff5ac3a7f95646f6123b0",
    importpath = "github.com/Masterminds/goutils",
)

# Indirect dependency
go_repository (
    name = "com_github_masterminds_semver",
    commit = "0fd41f6ff0825cf7efae00e706120bdd48914d93",
    importpath = "github.com/Masterminds/semver",
)

# Indirect dependency
go_repository(
    name = "com_github_imdario_mergo",
    commit = "4c317f2286be3bd0c4f1a0e622edc6398ec4656d",
    importpath = "github.com/imdario/mergo",
)

# Indirect dependency
go_repository(
    name = "com_github_google_uuid",
    commit = "c2e93f3ae59f2904160ceaab466009f965df46d6",
    importpath = "github.com/google/uuid",
)

# Indirect dependency
go_repository(
    name = "com_github_huandu_xstrings",
    commit = "8bbcf2f9ccb55755e748b7644164cd4bdce94c1d",
    importpath = "github.com/huandu/xstrings",
)

# Secure cookie storage
go_repository(
    name = "com_github_gorilla_securecookie",
    importpath = "github.com/gorilla/securecookie",
    sum = "h1:miw7JPhV+b/lAHSXz4qd/nN9jRiAFV5FwjeKyCS8BvQ=",
    version = "v1.1.1",
)

# Session utilities
go_repository(
    name = "com_github_gorilla_sessions",
    commit = "4355a998706e83fe1d71c31b07af94e34f68d74a",
    importpath = "github.com/gorilla/sessions",
)

# HTTP router
go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:gnP5JzjVOuiZD07fKKToCAOjS0yOpj/qPETTXCCS6hw=",
    version = "v1.7.3",
)

# Better logger
go_repository(
    name = "com_github_google_logger",
    commit = "7047ffcb7339f3f59be32de74a92217cb17cb40c",
    importpath = "github.com/google/logger",
)

# Syntax highlighter
go_repository(
    name = "com_github_alecthomas_chroma",
    commit = "f8432cf78f68e5adf203ad5cefaaf6244650b4d1",
    importpath = "github.com/alecthomas/chroma",
)

# Indirect dependency
go_repository(
    name = "com_github_danwakefield_fnmatch",
    importpath = "github.com/danwakefield/fnmatch",
    sum = "h1:y5HC9v93H5EPKqaS1UYVg1uYah5Xf51mBfIoWehClUQ=",
    version = "v0.0.0-20160403171240-cbb64ac3d964",
)

# Indirect dependency
go_repository(
    name = "com_github_dlclark_regexp2",
    importpath = "github.com/dlclark/regexp2",
    sum = "h1:CqB4MjHw0MFCDj+PHHjiESmHX+N7t0tJzKvC6M97BRg=",
    version = "v1.1.6",
)
