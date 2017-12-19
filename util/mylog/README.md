# mylog
日志模块

## Installation
```bash
glide get -v gitlab.matrix.netease.com/billing_dev/goal/pkg/mylog
```

## Usage and Examples

    package main

    import (
        "os"
        "path/filepath"
        "runtime"

        "gitlab.matrix.netease.com/billing_dev/goal/pkg/mylog"
        "gitlab.matrix.netease.com/billing_dev/goal/pkg/rotate"
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
