# podns
Run arbitrary cmd in pod namespace

# use-case
- When you don't want to add your tools to the pod but they're available on the host.
- When you just want to inspect a running process but don't feel like jumping through the latest hoops of the distributed container orchestration maze.
- When you need sudo power but the container doesn't run as root.

# requirements
- bash, oc cli and jq on local host
- docker, nsenter, ssh access to-, and sudoer on remote host

# usage

```bash
$ ./podns.sh <pod> <remote_user> [<cmd>]
```

`pod` and `remote_user` are required. `cmd` is optional - omit to get a remote shell on the host tied to the pod namespace. The pod should probably be in a running state. if `cmd` is set it is run with sudo. Use `@PID` to inject container pid in cmd. Make sure to put a long cmd in quotes.

# examples

List all network connections
```bash
# redirect stderr to ignore chatter about UID
$ ./podns.sh pod user "lsof -anP -i 2> /dev/null"
```

Poll every 5s for your server to be connected to redis
```bash
$ ./podns.sh pod user "lsof -anP -i:6379 -r 5"
```

List servers listening
```bash
$ ./podns.sh pod user "lsof -anP -i tcp" | grep -i listen
```

Check http end point
```bash
$ ./podns.sh pod user "curl localhost:5000/health -sik"
```

Get a shell
```bash
$ ./podns.sh pod user
```

# todos
- [ ] if $1 is --help print usage and exit
- [x] validate output
- [ ] put echo steps behind debug flag
- [x] shell option
- [x] arg REMOTE_USER
- [ ] add example on pipe to cmd that defines connection name from host:port
- [x] arg cmd
- [x] echo to stderr so pipes work
- [ ] add strace example
- [ ] allow multiple pids or dc
