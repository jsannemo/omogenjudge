#include <iostream>
#include <string>
#include <thread>
#include <vector>

#include "container/container.h"
#include "gflags/gflags.h"
#include "proto/omogenexec.pb.h"
#include "util/format.h"
#include "util/log.h"

using std::cout;
using std::endl;
using std::string;
using std::thread;
using std::vector;

DEFINE_double(cputime, 10, "The CPU time limit of the process in seconds");
DEFINE_double(walltime, -1, "The wall time limit of the process in seconds. Default is the CPU time + 1 second");
DEFINE_double(memory, 200, "The memory limit in megabytes (10^6 bytes)");
DEFINE_double(diskio, 1000, "The disk IO limit in megabytes (10^6 bytes)");
DEFINE_double(processes, 1, "The number of processes allowed");

DEFINE_string(in, "", "An input file that is used as stdin");
DEFINE_string(out, "", "An output file that is used as stdout");
DEFINE_string(err, "", "An error file that is used as stderr");

DEFINE_string(exec, "", "The file path to the executable file");

DEFINE_string(dirs, "", "Directory rules for mounting things inside the container. Format is oldpath:newpath[:writable] separated by semicolon, e.g. /newtmp:/tmp:writable;/home:home");

// TODO(jsannemo): write flag validators
int main(int argc, char** argv) {
    gflags::ParseCommandLineFlags(&argc, &argv, true);
    omogenexec::InitLogging(argv[0]);
    if (FLAGS_walltime == -1) {
        FLAGS_walltime = FLAGS_cputime + 1;
    }

    ExecuteRequest request;
    request.mutable_limits()->set_cputime(FLAGS_cputime);
    request.mutable_limits()->set_walltime(FLAGS_walltime);
    request.mutable_limits()->set_memory(FLAGS_memory * 1000);
    request.mutable_limits()->set_diskio(FLAGS_diskio * 1000);
    request.mutable_limits()->set_processes(FLAGS_processes);

    request.mutable_streams()->set_infile(FLAGS_in);
    request.mutable_streams()->set_outfile(FLAGS_out);
    request.mutable_streams()->set_errfile(FLAGS_err);

    vector<string> rules = omogenexec::Split(FLAGS_dirs, ';');
    for (const string& ruleStr : rules) {
        vector<string> fields = omogenexec::Split(ruleStr, ':');
        assert(2 <= fields.size() && fields.size() <= 3);
        DirectoryRule *rule = request.add_directories();
        rule->set_oldpath(fields[0]);
        rule->set_newpath(fields[1]);
        if (fields.size() == 3) {
            assert(fields[2] == "writable");
            rule->set_writable(true);
        }
    }

    request.mutable_command()->set_command(string(argv[1]));
    for (int i = 2; i < argc; ++i) {
        request.mutable_command()->add_flags(string(argv[i]));
    }

    omogenexec::Container c;
    cout << c.Execute(request).DebugString();
}
