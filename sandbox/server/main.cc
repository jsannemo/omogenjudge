#include "gflags/gflags.h"
#include "grpc++/grpc++.h"
#include "grpc++/security/server_credentials.h"
#include "grpc++/server.h"
#include "grpc++/server_builder.h"
#include "sandbox/server/executor.h"

using grpc::Server;
using grpc::ServerBuilder;

DEFINE_string(
    sandbox_listen_addr, "127.0.0.1:61810",
    "The address the sandbox server should listen to in the format host:port");

namespace omogen {
namespace sandbox {

void RunServer() {
  ExecuteServiceImpl service;
  ServerBuilder builder;
  // TODO(jsannemo): this should not use insecure credentials
  builder.AddListeningPort(FLAGS_sandbox_listen_addr,
                           grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<Server> server(builder.BuildAndStart());
  LOG(INFO) << "Server started on " << FLAGS_sandbox_listen_addr << std::endl;
  server->Wait();
}

}  // namespace sandbox
}  // namespace omogen

int main(int argc, char** argv) {
  gflags::SetUsageMessage("Run the OmogenJudge sandbox component");
  gflags::ParseCommandLineFlags(&argc, &argv, true);
  google::InitGoogleLogging(argv[0]);
  google::InstallFailureSignalHandler();
  omogen::sandbox::RunServer();
  gflags::ShutDownCommandLineFlags();
}
