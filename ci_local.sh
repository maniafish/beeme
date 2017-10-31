#!/bin/sh

sh -c "mysql -u root -e 'create database orm_test;'"
rm beeme
golint -set_exit_status $(go list ./... | grep -v /vendor/)
if [ $? -ne 0 ];then
    exit 1
fi

go vet $(go list ./...|grep -v vendor/)
if [ $? -ne 0 ];then
    exit 1
fi

sh build.sh
if [ $? -ne 0 ];then
    exit 1
fi

go test -v $(go list ./...|grep -v vendor/)
