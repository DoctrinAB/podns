#!/bin/bash

# Run arbitrary cmd in pod namespace
# version 0.4
# $ ./podns.sh <pod> <remote_user> [<cmd>]

# log ops to stderr so pipes work
function log() {
	# use cat <<< "$@" to allow for args like -n to be printed
	echo "$@" 1>&2;
}

POD="$1"
REMOTE_USER="$2"
CMD="$3"

# require at least 2 args
if test $# -lt 2; then
	log "error: missing args"
	exit 1
fi

log "* get node_ip and container_id from pod"
BATCH=`kubectl get pod $POD -o=jsonpath='{.status.hostIP},{.status.containerStatuses[0].containerID}'`
if test -z "$BATCH"; then
	log "error: no items found"
	exit 1
fi
NODE_IP=`cut -d , -f 1 <<< $BATCH`
log "$NODE_IP"
CONTAINER_ID=`cut -d , -f 2 <<< $BATCH | sed 's/docker:\/\///'`
log "$CONTAINER_ID"

log "* get pid"
PID=`ssh $REMOTE_USER@$NODE_IP "sudo docker inspect -f '{{.State.Pid}}' $CONTAINER_ID"`
if test -z "$PID"; then
	log "error: pid not found"
	exit 1
fi
log "$PID"

if test -z "$CMD"; then
	log "* running a shell in pid namespace"
else
	# replace @PID w $PID
	CMD=`sed "s/@PID/$PID/" <<< "$CMD"`
	log "* running $CMD in pid namespace"
fi

ssh $REMOTE_USER@$NODE_IP "sudo nsenter -t $PID -n $CMD"
