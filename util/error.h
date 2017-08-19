#include <string>

using namespace std;

[[noreturn]] void crash();
[[noreturn]] void crashSyscall(string file, int line, string syscall);
#define CRASH() crash()
#define CRASH_ERROR(syscall) crashSyscall(__FILE__, __LINE__, syscall)
