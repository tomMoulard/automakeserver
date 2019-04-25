package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    "encoding/json"
    "os"
    "os/exec"
    "github.com/tommoulard/automakeserver/helper"
)

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
        cmd := exec.Command("/bin/sh", "script.sh", request.Action)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()
    }(&request)
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/update/", Update)
    log.Fatal(http.ListenAndServe(":8080", router))
}
