package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var portToServe int

func init() {
	const (
		defaultPort = 8080
		usage       = "the port to serve"
	)
	flag.IntVar(&portToServe, "port", defaultPort, usage)
	flag.IntVar(&portToServe, "p", defaultPort, usage+" (shorthand)")
}

func main() {
	flag.Parse()

	str := `{"nodes":{"10.4.0.1":{"type":"Standby","agent":"UP","db":"UP","vip":"10.4.0.227","vip_active":false,"xlog":"3079\/64650BF8","xloginfo":""},"10.4.0.12":{"type":"Master","agent":"UP","db":"UP","vip":"10.4.0.227","vip_active":true,"xlog":"3079\/64650BF8","xloginfo":""},"10.4.0.121":{"type":"Witness","agent":"UP","db":"N\/A","vip":"10.4.0.227","vip_active":false},"10.4.0.122":{"type":"Witness","agent":"UP","db":"N\/A","vip":"10.4.0.227","vip_active":false}},"allowednodes":["10.4.0.122","10.4.0.121","10.4.0.12","10.4.0.1"],"membershipcoordinator":"10.4.0.122","failoverpriority":["10.4.0.1"],"minimumstandbys":0,"missingnodes":[],"messages":[]}`
	res := efm{}
	json.Unmarshal([]byte(str), &res)

	foo := newNodeCollector(res)
	prometheus.MustRegister(foo)

	http.Handle("/metrics", promhttp.Handler())
	log.Info(fmt.Sprintf("Beginning to serve on port %d", portToServe))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portToServe), nil))
}
