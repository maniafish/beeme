# rotate
日志轮转模块

## Installation

    go get github.com/maniafish/beeme/util/rotate

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
        filename := "test.log"
        logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
        if err != nil {
            panic(err)
        }

        // set logger
        _, lastFilePath, _, _ := runtime.Caller(0)
        ignoredDirPrefix := filepath.Dir(filepath.Dir(lastFilePath))
        mylog.Init("DEBUG", logFile, mylog.Lfile|mylog.Lrelative|mylog.Ljson|mylog.Lwelog, ignoredDirPrefix, os.Stdout)
        Log := mylog.GetStdLog()

        // predefine fields
        m := map[string]interface{}{
            "pid": os.Getpid(),
        }

        Log.Predefine(m)

        Log.Out = rotate.InitRotate(Log, filename, rotate.Daily, 30)

        // log output
        Log.Infof("test log")
    }
