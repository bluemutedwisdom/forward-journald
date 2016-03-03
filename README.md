# forward-journald

Forward stdin to journald

The main driver for this program is < go 1.6rc2 has a issue where 10 SIGPIPE's on stdout or stderr cause go to generate a
non-trappable SIGPIPE killing the process.  This happens when journald is restarted while docker is running under systemd.

Changes for `docker.service`:

```
:
[Service]
Type=notify
NotifyAccess=all
ExecStart=/bin/sh -c '/usr/bin/docker daemon \
          $OPTIONS \
          $DOCKER_STORAGE_OPTIONS \
          $DOCKER_NETWORK_OPTIONS \
          $INSECURE_REGISTRY \
          2>&1 | /usr/bin/forward-journald -tag docker'
StandardOutput=null
StandardError=null
:
```

## Bugs
https://bugzilla.redhat.com/show_bug.cgi?id=1300076

