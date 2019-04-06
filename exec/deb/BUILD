load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

pkg_tar(
    name = "exec_server",
    srcs = [
        "//exec/server:omogenexec_server",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/exec/server",
)

pkg_tar(
    name = "exec_cli",
    srcs = [
        "//exec/cli:omogenexec_cli",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/exec/cli",
)

pkg_tar(
    name = "config",
    srcs = [
        "//exec/config:users.conf",
    ],
    mode = "0644",
    package_dir = "/etc/omogen/exec",
    strip_prefix = "/exec/config",
)

pkg_tar(
    name = "initd",
    srcs = [
        "//exec/deb:omogenexecd",
    ],
    mode = "0755",
    package_dir = "/etc/init.d",
    strip_prefix = "/exec/deb",
)

pkg_tar(
    name = "exec_runner",
    srcs = [
        "//exec/container/inside:omogenexec_runner",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/exec/container/inside",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":config",
        ":exec_cli",
        ":exec_server",
        ":exec_runner",
        ":initd",
    ],
)

pkg_deb(
    name = "omogenexec",
    architecture = "amd64",
    built_using = "bazel",
    data = ":debian-data",
    description = "Sandbox service for untrusted applications",
    maintainer = "Johan Sannemo <johan.sannemo@kodsport.se>",
    package = "omogenexec",
    postinst = "debian/postinst",
    postrm = "debian/postrm",
    prerm = "debian/prerm",
    version = "0.0.1",
)
