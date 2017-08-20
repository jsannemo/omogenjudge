#include <algorithm>
#include <cassert>
#include <gflags/gflags.h>
#include <iostream>
#include <string>

#include "error.h"
#include "log.h"

using std::cerr;
using std::endl;
using std::ostream;
using std::streambuf;
using std::string;

// TODO(jsannemo): write a flag validator for this to avoid crashes
DEFINE_string(loglevel, "info", "The lowest log level to display (trace, debug, info, warn, error, fatal, none)");

namespace omogenexec {

static string toLevelString(const LogLevel level) {
    if (level == LogLevel::TRACE) return "trace";
    if (level == LogLevel::DEBUG) return "debug";
    if (level == LogLevel::INFO) return "info";
    if (level == LogLevel::WARN) return "warn";
    if (level == LogLevel::ERROR) return "error";
    if (level == LogLevel::FATAL) return "fatal";
    if (level == LogLevel::NONE) return "none";
    assert(false);
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
    OE_LOG(FATAL) << "Invalid log level " << level << endl;
	crash();
}

static LogLevel logThreshold;
static string loggerName;

void InitLogging(const string& name) {
    loggerName = name;
    logThreshold = toLogLevel(FLAGS_loglevel);
}

// When the log level threshold is too low, we wish to redirect any
// log levels to a dummy stream instead of cerr. This streambuf implementation
// is simply a no-op for printing, used for this purpose.
class NullBuffer : public streambuf {
public:
    int overflow(int c) { return c; }
};

static NullBuffer nullBuffer;
static ostream nullStream(&nullBuffer);

ostream& loggerForLevel(LogLevel level) {
    if (level < logThreshold) {
        return nullStream;
    }
    cerr << "[" << loggerName << "][" << toLevelString(level) << "]";
    return cerr;
}

}
