# The sandbox

The OmogenJudge sandbox is a custom built Linux sandbox using namespaces, containers and qouta to protect against malicious programs and measure resource usage.
The code for it lies in the directory `sandbox/`.

## The execution flow
The execution flow if a contained program consists of several steps.

1. An execution request is received by the gRPC server (`omogenjudge-sandbox`).
1. A container is allocated for the execution request.
1. The container starts up the container init process (`omogenjudge-sandboxr`).
1. An execution from the request is sent to the container.
1. The container forwards the request to the init process.
1. The init process is moved into one of the sandbox users, and quota for it is set.
1. The init process forks, set up an environment and execve()s's the requested command.
1. The init process waits for the program to finish, and reports its termination back to the container.
1. The container sends the results back.
1. If further executions come, steps 4-8 are repeated for each one.
1. The container is de-allocated.

A single allocated container can thus execute multiple programs.

If the executed program exceeds some resource limit, it is killed and a new container will be allocated if further execution requests are sent.

## Resource usage measurement
Resource usage (time and memory) are measured through the use of cgroups.
Wall-time is measured by a monitoring process in the container.
Disk usage is not measured.

## Containment mechanisms
There are three aspects of containment:

- preventing excessive resource usage
- preventing information leakage
- preventing interference with the rest of the system

Disk (blocks and inodes) is controlled using the standard Linux quota mechanisms.
Time is measured through cgroups, and if usage exceeds the allocated time, the monitor thread will kill the container init process.
Memory limits is enforced by the group.
Processes are restriced using rlimits -- cgroup counting of PIDs is very asynchronous, and can sometimes count processes that have already died.
This makes cgroups unsuitable for restricting processes for us, since we execute many processes in the same container.

Information leakage is mainly prevented by restricting file access through a chroot jail, only allowing a view of a restricted section of the file system.
The container has many new namespaces (for example networks) that restrict information leakage about existing users, network services, running processes.

Interference with the system is prevented by the use of several own namespaces, running the program as a low-privileged user, and restricting most parts of the file system as read-only.

## Sandbox users
When installing the sandbox package, it sets up a large number of users to be used for running the sandbox.
They all belong to the secondary group `omogenjudge-clients`, which can thus be used for setting file permissions without having world read-writable files.
It is also these users that quota is enforced for.

## The `omogenjudge-sandbox` user
The sandbox server is itself run with a low-privileged user: `omogenjudge-sandbox`.
This only has write access to the `/var/lib/omogen/sandbox` folder, which it uses to store the root filesystems of its containers.

## Reclaiming ownership
A sandboxed program may create files of its own.
This is a problem, since the unprivileged user process may not be allowed to run them.
Thus, users of the sandbox will likely need a root way of reclaiming ownership of created files.
In OmogenJudge, the evaluation component uses a `setuid` program
