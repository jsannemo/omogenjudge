cc_binary(
    name = "bad_programs_test",
    srcs = ["bad_programs_test.cc"],
    copts = ["-Iexternal/gtest/include"],
    deps = [
        "@gtest//:main",
        "//container:container"
    ],
    data = [":programs"]
)
