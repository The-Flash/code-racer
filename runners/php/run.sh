#!/bin/sh
maxProcesses=$1
maxOpenFiles=$2
maxFileSize=$3
timeoutSecs=$4
shift 4
prlimit --nproc=$maxProcesses --nofile=$maxOpenFiles --fsize=$maxFileSize timeout -s SIGKILL -k $timeoutSecs $timeoutSecs php $@