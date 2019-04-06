#include <gflags/gflags.h>
#include <grpc++/grpc++.h>
#include <grpc++/security/server_credentials.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>

#include "exec/server/executor.h"

using grpc::Server;
using grpc::ServerBuilder;

DEFINE_string(listen, "127.0.0.1:61810", "The address and port to listen to");

namespace omogen {
namespace exec {

void RunServer() {
  ExecuteServiceImpl service;
  ServerBuilder builder;
  builder.AddListeningPort(FLAGS_listen, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<Server> server(builder.BuildAndStart());
  LOG(INFO) << "Server started on " << FLAGS_listen << std::endl;
  server->Wait();
}

}  // namespace exec
}  // namespace omogen

int main(int argc, char** argv) {
  gflags::SetUsageMessage("Run an executor daemon");
  gflags::ParseCommandLineFlags(&argc, &argv, true);
  google::InitGoogleLogging(argv[0]);
  omogen::exec::RunServer();
  gflags::ShutDownCommandLineFlags();
}
