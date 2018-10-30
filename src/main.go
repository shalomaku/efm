package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var portToServe int
var commandToExecute string

func init() {
	const (
		defaultPort      = 8080
		defaultPortUsage = "The port to serve"

		defaultCommand      = "/usr/edb/efm-3.2/bin/efm cluster-status-json efm"
		defaultCommandUsage = "The command to execute"
	)
	flag.IntVar(&portToServe, "port", defaultPort, defaultPortUsage)
	flag.IntVar(&portToServe, "p", defaultPort, defaultPortUsage+" (shorthand)")

	flag.StringVar(&commandToExecute, "command", defaultCommand, defaultCommandUsage)
	flag.StringVar(&commandToExecute, "c", defaultCommand, defaultCommandUsage+" (shorthand)")
}

func executeCommand(cmdToExecute string) []byte {
	log.Info("Running command " + cmdToExecute)
	str := strings.Split(cmdToExecute, " ")
	out, _ := exec.Command(str[0], str[1:]...).Output()
	log.Info("Command was successfully executed.")

	return out
}

func convertOutputToEfm(str []byte) efm {
	res := efm{}
	log.Info("Beginning to decode output..")
	err := json.Unmarshal(str, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Json decode finished.")

	return res
}

func main() {
	flag.Parse()

	foo := newNodeCollector(convertOutputToEfm(executeCommand(commandToExecute)))
	prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Info(fmt.Sprintf("Beginning to serve on port %d", portToServe))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portToServe), nil))
}
