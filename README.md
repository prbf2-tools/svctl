# Svctl

Svctl operates as (Daemon)[#Daemon] which acts as process watcher and (CLI)[#CLI] which is an interface to communicate with Daemon.

## Daemon
Daemon keeps track of all started PRBF2 processes, stops them on user input and restarts if they exit on their own.

Daemon has to be running so other commands can interact with PRBF2 server instance.

### Windows
On Windows when a process crashes and raises an error, this error has to be closed manually,
Daemon is running an Error Killer which closes the process in such cases (and triggers restart mechanism).

### Linux
`/proc/sys/kernel/core_pattern`
```
"core.%e.%p" > /proc/sys/kernel/core_pattern
```

## Init

```bash
$ svctl init
```

## Register

```bash
$ svctl register
```
