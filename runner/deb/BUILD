load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

pkg_tar(
    name = "omogenrunner_server",
    srcs = [
        "//runner/server:omogenrunner_server",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/runner/server/linux_amd64_stripped",
)

pkg_tar(
    name = "omogenrunner_cli",
    srcs = [
        "//runner/cli:omogenrunner_cli",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/runner/cli/linux_amd64_stripped",
)

pkg_tar(
    name = "passwd",
    srcs = [
        "//runner/deb/etc:passwd",
    ],
    mode = "0644",
    package_dir = "/var/lib/omogen",
    strip_prefix = "/runner/deb",
)

pkg_tar(
    name = "initd",
    srcs = [
        "//runner/deb:omogenrunnerd",
    ],
    mode = "0755",
    package_dir = "/etc/init.d",
    strip_prefix = "/runner/deb",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":initd",
        ":omogenrunner_cli",
        ":omogenrunner_server",
        ":passwd",
    ],
)

pkg_deb(
    name = "omogenrunner",
    architecture = "amd64",
    built_using = "bazel",
    data = ":debian-data",
    description = "Program execution service",
    maintainer = "Johan Sannemo <johan.sannemo@kodsport.se>",
    package = "omogenrunner",
    postinst = "debian/postinst",
    prerm = "debian/prerm",
    version = "0.0.1",
)
