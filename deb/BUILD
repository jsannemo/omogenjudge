load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar", "pkg_deb")

pkg_tar(
    name = "omogenexec_daemon",
    strip_prefix = "/daemon",
    package_dir = "/usr/bin",
    srcs = [
				"//daemon:omogenexec_daemon",
		],
    mode = "0755",
)

pkg_tar(
    name = "omogenexec_cli",
    strip_prefix = "/cli",
    package_dir = "/usr/bin",
    srcs = [
				"//cli:omogenexec_cli",
		],
    mode = "0755",
)

pkg_tar(
    name = "config",
    strip_prefix = "/config",
    package_dir = "/etc/omogenexec",
    srcs = [
        "//config:users.conf",
    ],
    mode = "0644",
)

pkg_tar(
    name = "initd",
    strip_prefix = "/deb",
    package_dir = "/etc/init.d",
    srcs = [
        "//deb:omogenexecd",
    ],
    mode = "0755",
)

pkg_tar(
    name = "omogenexec_runner",
    strip_prefix = "/container/inside",
    package_dir = "/usr/bin",
    srcs = [
				"//container/inside:omogenexec_runner",
		],
    mode = "0755",
)

pkg_tar(
    name = "debian-data",
    extension = "tar.gz",
    deps = [
        ":config",
        ":initd",
        ":omogenexec_cli",
        ":omogenexec_daemon",
        ":omogenexec_runner",
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
    version = "0.0.1",
    postinst = "debian/postinst",
    prerm = "debian/prerm",
    postrm = "debian/postrm",
)

