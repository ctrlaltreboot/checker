package main

import (
	"encoding/json"
	"flag"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Checks []struct {
		Name string
		URL  string
	}
}

var cf Config

func init() {
	var conf string
	flag.StringVar(&conf, "config", "./config.json", "A fully qualified path for the configuration file")
	flag.Parse()

	f, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &cf)
	if err != nil {
		log.Fatal(err)
	}
}

func measure(address string) time.Duration {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	start := time.Now()
	_, err := netClient.Get(address)
	if err != nil {
		panic(err)
	}

	return time.Since(start)
}

func ddog(nameSpace string, resTime time.Duration) {
	cl, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}

	cl.Namespace = nameSpace

	err = cl.Timing("request.duration", resTime, nil, 1)
	if err != nil {
		log.Fatal(err)
	}
}

func check() {
	for _, c := range cf.Checks {
		r := measure(c.URL)
		ddog(c.Name, r)
		log.Println("Response time for ", c.URL, "is", r)
	}
}

func main() {
	s := gocron.NewScheduler()
	s.Every(20).Seconds().Do(check)
	<-s.Start()
}
