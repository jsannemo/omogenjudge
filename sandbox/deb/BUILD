load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_deb", "pkg_tar")

pkg_tar(
    name = "server",
    srcs = [
        "//sandbox/server:omogenjudge-sandbox",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/server",
)

pkg_tar(
    name = "config",
    srcs = [
        "//sandbox/config:sandbox.conf",
    ],
    mode = "0644",
    package_dir = "/etc/omogen/sandbox",
    strip_prefix = "/sandbox/config",
)

pkg_tar(
    name = "service",
    srcs = [
        "//sandbox/deb:omogenjudge-sandbox.service",
    ],
    mode = "0755",
    package_dir = "/etc/systemd/system",
    strip_prefix = "/sandbox/deb",
)

pkg_tar(
    name = "runner",
    srcs = [
        "//sandbox/container/inside:omogenjudge-sandboxr",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/container/inside",
)

pkg_tar(
    name = "clear",
    srcs = [
        "//sandbox/container/sandboxclear:omogenjudge-sandboxc",
    ],
    mode = "0755",
    package_dir = "/usr/bin",
    strip_prefix = "/sandbox/container/sandboxclear",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":config",
        ":service",
        ":runner",
        ":server",
        ":clear",
    ],
)

pkg_deb(
    name = "omogenjudge-sandbox",
    architecture = "amd64",
    built_using = "bazel",
    data = ":debian-data",
    description = "Sandbox service for untrusted applications",
    maintainer = "Johan Sannemo <johan.sannemo@kodsport.se>",
    package = "omogenjudge-sandbox",
    postinst = "debian/postinst",
    postrm = "debian/postrm",
    prerm = "debian/prerm",
    version = "0.0.1",
)
