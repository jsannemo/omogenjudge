cc_library(
        name = "container",
        visibility = ["//visibility:public"],
        hdrs = ["chroot.h", "cgroups.h", "container.h", "init.h"],
        srcs = ["chroot.cc", "cgroups.cc", "container.cc", "init.cc"],
        deps = [
                "//proto:omogenexec_cc_proto",
                "//util:util",
        ],
)
