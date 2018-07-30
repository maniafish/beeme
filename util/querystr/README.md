# querystr

   encoding structs into URL query parameters.
   source from github.com/google/go-querystring, add tag support

## Installation

    go get github.com/maniafish/beeme/util/querystr

## Usage and Examples

    package main

    import (
        "fmt"

        "github.com/maniafish/beeme/util/querystr"
    )

    func main() {
        type Param struct {
            GameID string `json:"gameid"`
            SN     string `json:"sn"`
        }

        p := &Param{"test", "test123"}
        v, err := querystr.Values(p, "json")
        if err != nil {
            panic(err)
        }

        // output: url.Values, map[gameid:[test] sn:[test123]]
        fmt.Printf("%T, %+v\n", v, v)
    }
