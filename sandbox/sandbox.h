#include "proto/omogenexec.pb.h"

/*
 * A process container, which can execute a command isolated from the rest of the system.
 */
class Container {

    pid_t pid;
    int commandPipe[2];
    int responsePipe[2];

    void killSupervisor();

public:
    Container();
    ~Container();

    ExecuteResponse Execute(const ExecuteRequest& request);

};
