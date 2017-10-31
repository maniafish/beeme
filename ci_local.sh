#!/bin/sh

sh -c "mysql -u root -e 'create database orm_test;'"
for pkg in $(go list ./... | grep -v /vendor/ | grep -v beeme) ; do golint -set_exit_status $pkg ; done
go vet $(go list ./...|grep -v vendor/)
sh build.sh
go test -v $(go list ./...|grep -v vendor/)
