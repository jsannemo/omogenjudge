# Installing omogenexec

## Configuring control groups
To use the control group features, a privileged user must create a parent control group called `omogencontain` that the user running execution has read and write rights to.

```bash
sudo mkdir /sys/fs/cgroup/{cpuacct,cpuset,pids,memory}/omogencontain
sudo chown containeruser:containeruser /sys/fs/cgroup/{cpuacct,cpuset,pids,memory}/omogencontain
```
