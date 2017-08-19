#include <cstdlib>
#include <cstring>
#include "log.h"

[[noreturn]] void crash() {
    exit(1);
}

[[noreturn]] void crashSyscall(string file, int line, string syscall) {
    LOG_LOCATION(FATAL, file, line) << syscall << ": " << strerror(errno) << endl;
    crash();
}
