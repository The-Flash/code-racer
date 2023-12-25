#!/bin/sh
maxProcesses=$1
maxOpenFiles=$2
maxFileSize=$3
timeoutSecs=$4
shift 4
filename=$1
prlimit --nproc=$maxProcesses --nofile=$maxOpenFiles --fsize=$maxFileSize timeout -s SIGKILL -k 10 10 java $filename $@
