#include <cstdlib>
#include <cstring>
#include "logger/logger.h"

void __attribute__((noreturn)) crash() {
    exit(1);
}

void __attribute__((noreturn)) crashSyscall(string syscall) {
    LOG(FATAL) << syscall << ": " << strerror(errno) << endl;
    crash();
}
