#include <string>
#include <ostream>

using namespace std;

enum LogLevel {
    TRACE = 0,
    DEBUG,
    INFO,
    WARN,
    ERROR,
    FATAL,
    NONE
};

namespace log {

void InitLogging(const string& name);

ostream& logAt(const LogLevel level);

}

#define LOG(level) LOG_LOCATION(level, __FILE__, __LINE__)
#define LOG_LOCATION(level, file, line) log::logAt(level) << "[" << file << ":" << line << "] "
