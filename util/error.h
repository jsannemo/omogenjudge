#pragma once

#include <string>

namespace omogenexec {

[[noreturn]] void crash();
[[noreturn]] void crashSyscall(std::string file, int line, std::string syscall);
#define OE_CRASH() omogenexec::crash()
#define OE_FATAL(syscall) omogenexec::crashSyscall(__FILE__, __LINE__, syscall)

}
