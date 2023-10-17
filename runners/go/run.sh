#!/bin/bash
filename=*.go
shift
go run $filename "$@"