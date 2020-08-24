# podns
Run arbitrary cmd in pod namespace

### use-case
- When you don't want to add your tools to the pod but they're available on the host.
- When you just want to inspect a running process but don't feel like jumping through the latest hoops of the distributed container orchestration maze.
- When you need sudo power but the container doesn't run as root.

### requirements
- bash and kubectl on local host
- docker, nsenter, ssh access to-, and sudoer on remote host
- Any command you want to run must be available on the remote host

# usage

```bash
$ ./podns.sh <pod> <remote_user> [<cmd>]
```

`pod` and `remote_user` are required. `cmd` is optional - omit to get a remote shell on the host tied to the pod namespace. The pod should probably be in a running state. if `cmd` is set it is run with sudo. Use `@PID` to inject container pid in cmd. Make sure to put a long cmd in quotes.

### examples

List all network connections
```bash
# redirect stderr to ignore chatter about UID
$ ./podns.sh pod user "lsof -anP -i 2> /dev/null"
```

Get a shell
```bash
$ ./podns.sh pod user
```

### useful commands

```bash
# Poll every 5s for your server to be connected to redis
lsof -anP -i:6379 -r 5
# List servers listening
lsof -anP -i tcp | grep -i listen
# Check http end point
curl localhost:5000/health -sik
# Run a trace on read syscalls
strace -e trace=read -t -s 128 -p @PID
# Grab mem usage
ps -v -p @PID
```

### filter
Use filters to filter output of a remote command by piping to it.

```bash
# get number of websockets connected to <port>
$ ./podns.sh pod user "lsof -nP -i 2> /dev/null" | go run filter/ws.go <port>
```

### kubectl plugin
podns can be used as a kubectl plugin and invoked with `kubectl podns ...args`. Just copy the script to somewhere in your $PATH and rename it to `kubectl-podns`.

```bash
$ cp -iv podns.sh /usr/local/bin/kubectl-podns
```

### multiple pods (beta)
Experimental support for running remote command in multiple pod ns.

```bash
$ go run podns.go remoteCmd pod [pod...] [-h help -u remoteUser]
```

# todos
- [ ] if $1 is --help print usage and exit
- [x] validate output
- [ ] put echo steps behind debug flag
- [x] shell option
- [x] arg REMOTE_USER
- [x] filters
- [x] arg cmd
- [x] echo to stderr so pipes work
- [x] add strace example
- [x] allow multiple pids or dc
