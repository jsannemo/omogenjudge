cc_library(
    name = "container",
    visibility = ["//visibility:public"],
    hdrs = ["container.h"],
    srcs = ["container.cc"],
    deps = [
            ":container-internals",
            "//proto:omogenexec_cc_proto",
            "//util:util",
    ],
)

cc_library(
    name = "container-internals",
    visibility = ["//container:__pkg__"],
    hdrs = ["chroot.h", "cgroups.h", "init.h"],
    srcs = ["chroot.cc", "cgroups.cc", "init.cc"],
    deps = [
            "//proto:omogenexec_cc_proto",
            "//util:util",
    ],
)
