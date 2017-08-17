#include "container.h"
#include "proto/omogenexec.pb.h"

int main(int argc, char** argv) {
    Container cont;
    ExecuteRequest er;

    er.mutable_command()->set_command("/usr/bin/g++");
    er.mutable_command()->add_flags("/compile/A.cpp");
    er.mutable_command()->add_flags("-o");
    er.mutable_command()->add_flags("/compile/a.out");

    er.mutable_limits()->set_cputime(4);
    er.mutable_limits()->set_walltime(5);
    er.mutable_limits()->set_memory(1000 * 1000);
    er.mutable_limits()->set_diskio(1000 * 1000);

    er.mutable_streams()->set_infile("/usr/include/ftw.h");
    er.mutable_streams()->set_outfile("/compile/derp");
    er.mutable_streams()->set_errfile("/compile/err");

    DirectoryRule *rule = er.add_directories();
    rule->set_source("/tmp/derp");
    rule->set_destination("/compile");
    rule->set_writable("true");


    cont.Execute(er);
}
