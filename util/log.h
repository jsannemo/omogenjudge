#pragma once

#include <string>
#include <ostream>

enum LogLevel {
    TRACE = 0,
    DEBUG,
    INFO,
    WARN,
    ERROR,
    FATAL,
    NONE
};

namespace omogenexec {

void InitLogging(const std::string& name);

std::ostream& loggerForLevel(const LogLevel level);

#define OE_LOG(level) OE_LOG_LOCATION(level, __FILE__, __LINE__)
#define OE_LOG_LOCATION(level, file, line) omogenexec::loggerForLevel(level) << "[" << file << ":" << line << "] "

}
