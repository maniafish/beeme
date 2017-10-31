#!/bin/sh

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
