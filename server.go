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

var (
    script_pmacct = ""
    script_snmp = ""
    errorLog *os.File
    logger *log.Logger
)

func UseRequest(request *structs.GitStruct) {
    actions := map[string] bool {
        "created": true,
        "deleted": true,
        "edited": false,
        "prereleased": false,
        "published": true,
        "unpublished": false,
    }
    repos := map[string] string {
        "pmacct": script_pmacct,
        "snmp-streamer": script_snmp,
    }
    if !actions[request.Action] {return}
    var repoName = request.Repository.FullName
    var script = repos[repoName]
    logger.Printf("POST request: '%s' for '%s'\n", request.Action, repoName)
    cmd := exec.Command("/bin/sh", script, request.Action)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
    logger.Println("GET  request")
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
    go UseRequest(&request)
}

func initLog(logFile string) {
    errorLog, err := os.OpenFile(logFile,
        os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("Could not open the logging file: %v\n", err)
        os.Exit(1)
    }
    logger = log.New(errorLog, "", log.LstdFlags)
}

func main() {
    parser := argparse.NewParser("automakeserver",
        "Get POST request from Github and launch scripts")
    port := parser.String("p", "port",
        &argparse.Options{Default: "8080", Help:"Server port"})
    scrpt_snmp := parser.String("n", "script_snmp",
        &argparse.Options{Default: "script.sh",
            Help:"script to execute for the snmp repository"})
    scrpt_pmacct := parser.String("m", "script_pmacct",
        &argparse.Options{Default: "script.sh",
            Help:"script to execute for the pmacct repository"})
    logFile := parser.String("l", "log",
        &argparse.Options{Default: "log_file.log", Help:"Logging file"})

    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        os.Exit(1)
    }
    script_snmp = *scrpt_snmp
    script_pmacct = *scrpt_pmacct
    var portBuilder strings.Builder
    portBuilder.WriteString(":")
    portBuilder.WriteString(*port)
    initLog(*logFile)

    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/update/", Update)
    log.Fatal(http.ListenAndServe(portBuilder.String(), router))
}
