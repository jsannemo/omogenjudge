# Quota
To limit disk access, the standard Linux quota system is used.
This is generally not enabled on new systems, so it must first be enabled for the devices that submitted programs can write to.
In the operations of the judge, they will only ever write to files in the `/var/lib/omogen/submissions` directory, so quota must be enabled for the device this directory is mounted at.

Quota first needs to be enabled on the correct device by editing `/etc/fstab`, finding the line corresponding to the correct device, and adding `usrquota` to the options.
After adding the option, remount the device.

Then, quota needs to be enabled on the system with the commands

```bash
sudo quotaon -v /
sudo quotacheck -ugm /
```

If the `/var/lib/omogen/submissions` directory is not mounted on the root device, be sure to replace `/` in the commands with the correct mount point.
