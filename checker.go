package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

var conf string

func init() {
	flag.StringVar(&conf, "config", "./config.json", "A fully qualified path for the configuration file")
	flag.Parse()
}

func loadConfig() Config {
	f, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err)
	}

	var cf Config
	err = json.Unmarshal(f, &cf)
	if err != nil {
		log.Fatal(err)
	}

	return cf
}

func measure(address string) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	start := time.Now()
	_, err := netClient.Get(address)
	if err != nil {
		panic(err)
	}

	log.Println("Load Time:", time.Since(start))
}

func check() {
	cf := loadConfig()
	for _, c := range cf.Checks {
		fmt.Println("Response time for:", c.URL)
		measure(c.URL)
	}
}

func main() {
	s := gocron.NewScheduler()
	s.Every(20).Seconds().Do(check)
	<-s.Start()
}
