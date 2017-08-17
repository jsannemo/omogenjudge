cc_binary(
        name = "omogencontain",
        deps = [":omogencontainlib"],
        srcs = ["main.cc"]
)

cc_library(
        name = "omogencontainlib",
        visibility = ["//visibility:public"],
        hdrs = ["chroot.h", "cgroups.h", "container.h", "init.h"],
        srcs = ["chroot.cc", "cgroups.cc", "container.cc", "init.cc"],
        deps = [
                "//errors:errors",
                "//logger:logger",
                "//proto:omogenexec_cc_proto",
                "//external:gflags",
        ],
)
