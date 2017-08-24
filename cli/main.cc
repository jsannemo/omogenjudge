#include <iostream>
#include <string>
#include <thread>
#include <vector>

#include "container/container.h"
#include "gflags/gflags.h"
#include "proto/omogenexec.pb.h"
#include "util/error.h"
#include "util/files.h"
#include "util/format.h"
#include "util/log.h"

using std::cout;
using std::endl;
using std::string;
using std::thread;
using std::vector;

static bool validateLimit(const char* flagname, double value) {
    if (value < 0) {
        OE_LOG(FATAL) << "Invalid value for -" << flagname << ": " << value << endl;
        return false;
    }
    return true;
}

static bool validateLimitWithDefault(const char* flagname, double value) {
    if (value < 0 && value != -1) {
        OE_LOG(FATAL) << "Invalid value for -" << flagname << ": " << value << endl;
        return false;
    }
    return true;
}

static bool validateNonEmptyString(const char* flagname, const string& value) {
    if (value.empty()) {
        OE_LOG(FATAL) << "Invalid value for -" << flagname << ": " << value << endl;
        return false;
    }
    return true;
}

static bool validateDirectoryRules(const char *flagname, const string& value) {
    vector<string> rules = omogenexec::Split(value, ';');
    for (const string& rule : rules) {
        vector<string> fields = omogenexec::Split(rule, ':');
        if (fields.size() < 2 || 3 < fields.size()) {
            OE_LOG(FATAL) << "Directory rule " << rule << " invalid" << endl;
            return false;
        }
        if (fields[0].empty() || fields[0][0] != '/') {
            OE_LOG(FATAL) << "Path " << fields[0] << " is not absolute" << endl;
            return false;
        }
        if (fields[1].empty() || fields[1][0] != '/') {
            OE_LOG(FATAL) << "Path " << fields[1] << " is not absolute" << endl;
            return false;
        }
        if (!omogenexec::DirectoryExists(fields[0])) {
            OE_LOG(FATAL) << "Directory outside sandbox " << fields[0] << " does not exist" << endl;
            return false;
        }
        if (fields.size() == 3 && fields[2] != "writable") {
            OE_LOG(FATAL) << "Invalid directory options " << fields[2] << endl;
            return false;
        }
    }
    return true;
}

DEFINE_double(cputime, 10, "The CPU time limit of the process in seconds");
DEFINE_validator(cputime,  &validateLimit);
DEFINE_double(walltime, -1, "The wall time limit of the process in seconds. Default is the CPU time + 1 second");
DEFINE_validator(walltime,  &validateLimitWithDefault);
DEFINE_double(memory, 200, "The memory limit in megabytes (10^6 bytes)");
DEFINE_validator(memory,  &validateLimit);
DEFINE_double(diskio, 1000, "The disk IO limit in megabytes (10^6 bytes)");
DEFINE_validator(diskio,  &validateLimit);
DEFINE_double(processes, 1, "The number of processes allowed");
DEFINE_validator(processes,  &validateLimit);

DEFINE_string(in, "", "An input file that is used as stdin");
DEFINE_validator(in, &validateNonEmptyString);
DEFINE_string(out, "", "An output file that is used as stdout");
DEFINE_validator(out, &validateNonEmptyString);
DEFINE_string(err, "", "An error file that is used as stderr");
DEFINE_validator(err, &validateNonEmptyString);

DEFINE_string(dirs, "", "Directory rules for mounting things inside the container. Format is oldpath:newpath[:writable] separated by semicolon, e.g. /newtmp:/tmp:writable;/home:home");
DEFINE_validator(dirs, &validateDirectoryRules);

void parseDirectoryRule(DirectoryRule* rule, const string& ruleStr) {
    vector<string> fields = omogenexec::Split(ruleStr, ':');
    assert(2 <= fields.size() && fields.size() <= 3);
    rule->set_oldpath(fields[0]);
    rule->set_newpath(fields[1]);
    if (fields.size() == 3) {
        assert(fields[2] == "writable");
        rule->set_writable(true);
    }
}

void parseDirectoryRules(ExecuteRequest& request, const string& rulesStr) {
    vector<string> rules = omogenexec::Split(rulesStr, ';');
    for (const string& ruleStr : rules) {
        DirectoryRule *rule = request.add_directories();
        parseDirectoryRule(rule, ruleStr);
    }
}

void setLimits(ResourceLimits* limits) {
    limits->set_cputime(FLAGS_cputime);
    limits->set_walltime(FLAGS_walltime);
    limits->set_memory(FLAGS_memory * 1000);
    limits->set_diskio(FLAGS_diskio * 1000);
    limits->set_processes(FLAGS_processes);
}

void setStreams(StreamRedirections* streams) {
    streams->set_infile(FLAGS_in);
    streams->set_outfile(FLAGS_out);
    streams->set_errfile(FLAGS_err);
}

int main(int argc, char** argv) {
    gflags::ParseCommandLineFlags(&argc, &argv, true);
    omogenexec::InitLogging(argv[0]);
    if (FLAGS_walltime == -1) {
        FLAGS_walltime = FLAGS_cputime + 1;
    }

    ExecuteRequest request;
    setLimits(request.mutable_limits());
    setStreams(request.mutable_streams());
    parseDirectoryRules(request, FLAGS_dirs);

    if (argc < 2) {
        OE_LOG(FATAL) << "No command provided" << endl;
        OE_CRASH();
    }
    request.mutable_command()->set_command(string(argv[1]));
    for (int i = 2; i < argc; ++i) {
        request.mutable_command()->add_flags(string(argv[i]));
    }

    omogenexec::Container c;
    cout << c.Execute(request).DebugString();
}
