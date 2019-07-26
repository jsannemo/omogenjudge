#include <string>

#include "api/omogenexec.pb.h"
#include "container/outside/container.h"
#include "glog/logging.h"
#include "gtest/gtest.h"
#include "util/files.h"

using std::endl;
using std::ifstream;
using std::string;

using namespace omogenexec;
using namespace omogenexec::api;

void setRwDir(DirectoryMount* dr, const string& path) {
  dr->set_path_outside_container(path);
  dr->set_path_inside_container(path);
  dr->set_writable(true);
}

void setRDir(DirectoryMount* dr, const string& path) {
  dr->set_path_outside_container(path);
  dr->set_path_inside_container(path);
  dr->set_writable(false);
}

void addLimit(ResourceAmount* limit, ResourceType type, long long amount) {
  limit->set_type(type);
  limit->set_amount(amount);
}

void addStream(Streams::Mapping* mapping, const string& path, bool write) {
  mapping->set_write(write);
  Streams::Mapping::File* file = mapping->mutable_file();
  file->set_path_inside_container(path);
}

Termination compile_and_run(const string& path) {
  string tmp = MakeTempDir();
  WriteToFile(tmp + "/in", "");

  Execution compiler;
  compiler.mutable_command()->set_command("/usr/bin/g++");
  string programDir = string(get_current_dir_name()) +
                      "/container/tests/bad_programs/programs/";
  compiler.mutable_command()->add_flags(programDir + "/" + path);
  compiler.mutable_command()->add_flags("-o");
  compiler.mutable_command()->add_flags(tmp + "/a.out");
  addStream(compiler.mutable_environment()
                ->mutable_stream_redirections()
                ->add_mappings(),
            tmp + "/in", false);
  addStream(compiler.mutable_environment()
                ->mutable_stream_redirections()
                ->add_mappings(),
            tmp + "/out", true);
  addStream(compiler.mutable_environment()
                ->mutable_stream_redirections()
                ->add_mappings(),
            tmp + "/err", true);
  setRwDir(compiler.mutable_environment()->add_mounts(), tmp);
  setRDir(compiler.mutable_environment()->add_mounts(), programDir);
  addLimit(compiler.mutable_resource_limits()->add_amounts(),
           ResourceType::CPU_TIME, 2000);
  addLimit(compiler.mutable_resource_limits()->add_amounts(),
           ResourceType::WALL_TIME, 4000);
  addLimit(compiler.mutable_resource_limits()->add_amounts(),
           ResourceType::PROCESSES, 4);
  addLimit(compiler.mutable_resource_limits()->add_amounts(),
           ResourceType::MEMORY, 300 * 1000);

  Container compilerContainer;
  StatusOr<Termination> resp = compilerContainer.Execute(compiler);

  Execution exec;
  exec.mutable_command()->set_command(tmp + "/a.out");
  exec.mutable_command()->add_flags(tmp);
  exec.mutable_command()->add_flags(tmp);
  exec.mutable_command()->add_flags(tmp);
  addStream(
      exec.mutable_environment()->mutable_stream_redirections()->add_mappings(),
      tmp + "/in", false);
  addStream(
      exec.mutable_environment()->mutable_stream_redirections()->add_mappings(),
      tmp + "/out", true);
  addStream(
      exec.mutable_environment()->mutable_stream_redirections()->add_mappings(),
      tmp + "/err", true);
  setRwDir(exec.mutable_environment()->add_mounts(), tmp);
  addLimit(exec.mutable_resource_limits()->add_amounts(),
           ResourceType::CPU_TIME, 2000);
  addLimit(exec.mutable_resource_limits()->add_amounts(),
           ResourceType::WALL_TIME, 4000);
  addLimit(exec.mutable_resource_limits()->add_amounts(),
           ResourceType::PROCESSES, 20);
  addLimit(exec.mutable_resource_limits()->add_amounts(), ResourceType::MEMORY,
           300 * 1000);

  Container execContainer;
  StatusOr<Termination> resp2 = execContainer.Execute(exec);
  LOG(INFO) << tmp << endl;
  CHECK(resp2.ok()) << "Execution failed with " << resp2.status().error_code()
                    << ": " << resp2.status().error_message() << endl;
  // RemoveTree(tmp);
  return resp2.value();
}

TEST(BadPrograms, Brk) {
  Termination resp = compile_and_run("brk.cc");
  LOG(INFO) << "brk" << endl << resp.DebugString();
}

TEST(BadPrograms, Busy) {
  Termination resp = compile_and_run("busy.cc");
  LOG(INFO) << "busy" << endl << resp.DebugString();
}

TEST(BadPrograms, ForkBomb) {
  Termination resp = compile_and_run("forkbomb.cc");
  LOG(INFO) << "forkbomb" << endl << resp.DebugString();
}

TEST(BadPrograms, Malloc) {
  Termination resp = compile_and_run("malloc.cc");
  LOG(INFO) << "malloc" << endl << resp.DebugString();
}

TEST(BadPrograms, Vector) {
  Termination resp = compile_and_run("vector.cc");
  LOG(INFO) << "vector" << endl << resp.DebugString();
}
