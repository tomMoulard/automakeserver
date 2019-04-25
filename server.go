package main

import (
    "encoding/json"
    "fmt"
    "github.com/akamensky/argparse"
    "github.com/julienschmidt/httprouter"
    "github.com/tommoulard/automakeserver/helper"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
)

var script = ""

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")

    decoder := json.NewDecoder(r.Body)
    var request structs.GitStruct
    err := decoder.Decode(&request)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s", err)
        return
    }
    defer r.Body.Close()
    go func(request *structs.GitStruct) {
        actions := map[string] bool {
            "created": true,
            "deleted": true,
            "edited": false,
            "prereleased": false,
            "published": true,
            "unpublished": false,
        }
        if !actions[request.Action] {return}
        cmd := exec.Command("/bin/sh", script, request.Action)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()
    }(&request)
}

func main() {
    parser := argparse.NewParser("automakeserver","Get POST request from Github and launch scripts")
    port := parser.String("p", "port",
        &argparse.Options{Default: "8080", Help:"Server port"})
    scrpt := parser.String("s", "scipt",
        &argparse.Options{Default: "script.sh", Help:"Script to execute"})
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    script = *scrpt
    var portBuilder strings.Builder
    portBuilder.WriteString(":")
    portBuilder.WriteString(*port)

    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/update/", Update)
    log.Fatal(http.ListenAndServe(portBuilder.String(), router))
}
