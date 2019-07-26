load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

pkg_tar(
    name = "sandbox_server",
    srcs = [
        "//sandbox/server:omogensandbox_server",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/server",
)

pkg_tar(
    name = "sandbox_cli",
    srcs = [
        "//sandbox/cli:omogensandbox_cli",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/cli",
)

pkg_tar(
    name = "config",
    srcs = [
        "//sandbox/config:users.conf",
    ],
    mode = "0644",
    package_dir = "/etc/omogen/sandbox",
    strip_prefix = "/sandbox/config",
)

pkg_tar(
    name = "initd",
    srcs = [
        "//sandbox/deb:omogensandboxd",
    ],
    mode = "0755",
    package_dir = "/etc/init.d",
    strip_prefix = "/sandbox/deb",
)

pkg_tar(
    name = "sandbox_runner",
    srcs = [
        "//sandbox/container/inside:omogensandbox_runner",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/container/inside",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":config",
        ":initd",
        ":sandbox_cli",
        ":sandbox_runner",
        ":sandbox_server",
    ],
)

pkg_deb(
    name = "omogensandbox",
    architecture = "amd64",
    built_using = "bazel",
    data = ":debian-data",
    description = "Sandbox service for untrusted applications",
    maintainer = "Johan Sannemo <johan.sannemo@kodsport.se>",
    package = "omogensandbox",
    postinst = "debian/postinst",
    postrm = "debian/postrm",
    prerm = "debian/prerm",
    version = "0.0.1",
)
