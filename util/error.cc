#include <cstdlib>
#include <cstring>
#include <iostream>

#include "log.h"

using std::string;

namespace omogenexec {

[[noreturn]] void crash() {
    throw 1;
}

[[noreturn]] void crashSyscall(string file, int line, string syscall) {
    OE_LOG_LOCATION(FATAL, file, line) << syscall << ": " << strerror(errno) << std::endl;
    crash();
}

string StrError() {
    return string(strerror(errno));
}

}
