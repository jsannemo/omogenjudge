#include <string>

#include "gtest/gtest.h"
#include "container/container.h"
#include "proto/omogenexec.pb.h"
#include "util/files.h"
#include "util/log.h"

using std::endl;
using std::ifstream;
using std::string;

using namespace omogenexec;

void setRwDir(DirectoryRule* dr, const string& path) {
    dr->set_oldpath(path);
    dr->set_newpath(path);
    dr->set_writable(true);
}

void setRDir(DirectoryRule* dr, const string& path) {
    dr->set_oldpath(path);
    dr->set_newpath(path);
}

ExecuteResponse compile_and_run(const string& path) {
    string tmp = MakeTempDir();
    WriteToFile(tmp + "/in", "");

    ExecuteRequest compiler;
    compiler.mutable_command()->set_command("/usr/bin/g++");
    string programDir = string(get_current_dir_name()) + "/container/tests/bad_programs/programs/";
    compiler.mutable_command()->add_flags(programDir + "/" + path);
    compiler.mutable_command()->add_flags("-o");
    compiler.mutable_command()->add_flags(tmp + "/a.out");
    compiler.mutable_streams()->set_infile(tmp + "/in");
    compiler.mutable_streams()->set_errfile(tmp + "/err");
    compiler.mutable_streams()->set_outfile(tmp + "/out");
    DirectoryRule* dr = compiler.add_directories();
    setRwDir(dr, tmp);
    dr = compiler.add_directories();
    setRDir(dr, programDir);
    compiler.mutable_limits()->set_cputime(2);
    compiler.mutable_limits()->set_walltime(4);
    compiler.mutable_limits()->set_processes(4);
    compiler.mutable_limits()->set_memory(300 * 1000);
    compiler.mutable_limits()->set_diskio(300 * 1000);

    Container compilerContainer;
    ExecuteResponse resp = compilerContainer.Execute(compiler);

    ExecuteRequest exec;
    exec.mutable_command()->set_command(tmp + "/a.out");
    exec.mutable_command()->add_flags(tmp);
    exec.mutable_command()->add_flags(tmp);
    exec.mutable_command()->add_flags(tmp);
    exec.mutable_streams()->set_infile(tmp + "/in");
    exec.mutable_streams()->set_errfile(tmp + "/err");
    exec.mutable_streams()->set_outfile(tmp + "/out");
    dr = exec.add_directories();
    setRwDir(dr, tmp);
    exec.mutable_limits()->set_cputime(2);
    exec.mutable_limits()->set_walltime(4);
    exec.mutable_limits()->set_processes(20);
    exec.mutable_limits()->set_memory(300 * 1000);
    exec.mutable_limits()->set_diskio(300 * 1000);

    Container execContainer;
    ExecuteResponse resp2 = execContainer.Execute(exec);
    OE_LOG(INFO) << tmp << endl;
    //RemoveTree(tmp);
    return resp2;
}

TEST(BadPrograms, Brk) {
    ExecuteResponse resp = compile_and_run("brk.cc");
    OE_LOG(INFO) << "brk" << endl << resp.DebugString();
}

TEST(BadPrograms, Busy) {
    ExecuteResponse resp = compile_and_run("busy.cc");
    OE_LOG(INFO) << "busy" << endl << resp.DebugString();
}

TEST(BadPrograms, ForkBomb) {
    ExecuteResponse resp = compile_and_run("forkbomb.cc");
    OE_LOG(INFO) << "forkbomb" << endl << resp.DebugString();
}

TEST(BadPrograms, Malloc) {
    ExecuteResponse resp = compile_and_run("malloc.cc");
    OE_LOG(INFO) << "malloc" << endl << resp.DebugString();
}

TEST(BadPrograms, Vector) {
    ExecuteResponse resp = compile_and_run("vector.cc");
    OE_LOG(INFO) << "vector" << endl << resp.DebugString();
}

TEST(BadPrograms, WriteBomb2) {
    ExecuteResponse resp = compile_and_run("writebomb2.cc");
    OE_LOG(INFO) << "writebomb2" << endl << resp.DebugString();
}
