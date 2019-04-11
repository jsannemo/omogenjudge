load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

pkg_tar(
    name = "omogenrun_server",
    srcs = [
        "//run/server:omogenrun_server",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/run/server/linux_amd64_stripped",
)

pkg_tar(
    name = "omogenrun_cli",
    srcs = [
        "//run/cli:omogenrun_cli",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/run/cli/linux_amd64_stripped",
)

pkg_tar(
    name = "passwd",
    srcs = [
        "//run/deb/etc:passwd",
    ],
    mode = "0644",
    package_dir = "/var/lib/omogen",
    strip_prefix = "/run/deb",
)

pkg_tar(
    name = "initd",
    srcs = [
        "//run/deb:omogenrund",
    ],
    mode = "0755",
    package_dir = "/etc/init.d",
    strip_prefix = "/run/deb",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":initd",
        ":passwd",
        ":omogenrun_cli",
        ":omogenrun_server",
    ],
)

pkg_deb(
    name = "omogenrun",
    architecture = "amd64",
    built_using = "bazel",
    data = ":debian-data",
    description = "Program execution service",
    maintainer = "Johan Sannemo <johan.sannemo@kodsport.se>",
    package = "omogenrun",
    postinst = "debian/postinst",
    prerm = "debian/prerm",
    version = "0.0.1",
)
