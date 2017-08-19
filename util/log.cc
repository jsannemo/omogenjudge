#include <algorithm>
#include <gflags/gflags.h>
#include <iostream>
#include <string>

using namespace std;

#include "error.h"
#include "log.h"

DEFINE_string(loglevel, "info", "the lowest log level to display (trace, debug, info, warn, error, fatal, disabled)");

namespace log {

static string toLevelString(const LogLevel level) {
    if (level == LogLevel::TRACE) return "trace";
    if (level == LogLevel::DEBUG) return "debug";
    if (level == LogLevel::INFO) return "info";
    if (level == LogLevel::WARN) return "warn";
    if (level == LogLevel::ERROR) return "error";
    if (level == LogLevel::FATAL) return "fatal";
    if (level == LogLevel::NONE) return "none";
    LOG(FATAL) << "Invalid log level " << level << endl;
	crash();
}

static LogLevel toLogLevel(string level) {
    transform(level.begin(), level.end(), level.begin(),
            [](char c) { return (char)tolower(c); });
    if (level == "trace") return LogLevel::TRACE;
    if (level == "debug") return LogLevel::DEBUG;
    if (level == "info") return LogLevel::INFO;
    if (level == "warn") return LogLevel::WARN;
    if (level == "error") return LogLevel::ERROR;
    if (level == "fatal") return LogLevel::FATAL;
    if (level == "none") return LogLevel::NONE;
    LOG(FATAL) << "Invalid log level " << level << endl;
	crash();
}

static LogLevel leastLogLevel;
static string loggerName;

void InitLogging(const string& name) {
    loggerName = name;
    leastLogLevel = toLogLevel(FLAGS_loglevel);
}

class NullBuffer : public std::streambuf {
public:
    int overflow(int c) { return c; }
};

NullBuffer nullBuffer;
ostream nullStream(&nullBuffer);

ostream& logAt(const LogLevel level) {
    if (level < leastLogLevel) {
        return nullStream;
    }
    cerr << "[" << loggerName << "][" << toLevelString(level) << "]";
    return cerr;
}

}

