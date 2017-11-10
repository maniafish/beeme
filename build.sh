#!/bin/sh

function go_lint(){
    for pkg in $(go list ./... | grep -v /vendor/)
    do
        golint -set_exit_status $pkg
        if [ $? -ne 0 ]; then
            exit 1
        fi
    done
}

if [ "$1" == "lint" ]; then
    go_lint
    exit 0
fi

go build -v
if [ $? -ne 0 ];then
    exit 1
fi

# try start beeme for build router
./beeme &
sleep 1
PID=`ps ax | awk '$NF=="./beeme"{print $1; exit}'`
if [ -z "PID" ];then
    echo "start beeme failed"
    exit 1
fi

# kill proc
kill ${PID}
