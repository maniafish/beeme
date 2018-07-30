[![Build Status](https://travis-ci.org/maniafish/beeme.svg?branch=master)](https://travis-ci.org/maniafish/beeme)
# beeme
a web api by beego

## develop workflow

1. make run
2. view api doc at "http://127.0.0.1:8080/swagger/"
3. call interface by "127.0.0.1:8080/v1/..."

## release workflow

1. Every pull request must be named by $version-$msg
2. Add [CHANGELOG](./CHANGELOG.md) with pull request of this version
3. commit CHANGELOG
4. make && ./beeme-$branch-$version
