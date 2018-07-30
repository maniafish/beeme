# mylog
日志模块

## Installation
```bash
go get github.com/maniafish/beeme/util/mylog
```

## Usage and Examples

    package main

    import (
        "os"
        "path/filepath"
        "runtime"

        "github.com/maniafish/beeme/util/mylog"
        "github.com/maniafish/beeme/util/rotate"
    )

    func main() {
        // set logger
        _, lastFilePath, _, _ := runtime.Caller(0)
        ignoredDirPrefix := filepath.Dir(filepath.Dir(lastFilePath))
        mylog.Init("DEBUG", os.Stderr, mylog.Lfile|mylog.Lrelative|mylog.Ljson|mylog.Lwelog, ignoredDirPrefix, os.Stdout)
        Log := mylog.GetStdLog()

        // predefine fields
        m := map[string]interface{}{
            "pid": os.Getpid(),
        }

        Log.Predefine(m)

        // log output
        Log.Infof("test log")
    }
