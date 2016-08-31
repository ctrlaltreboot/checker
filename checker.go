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

type Urls struct {
	Checks []string
}

func loadConfig(path *string) Urls {
	f, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	var u Urls
	err = json.Unmarshal(f, &u)
	if err != nil {
		log.Fatal(err)
	}

	return u
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
	confPtr := flag.String("config", "./config.json", "A fully qualified path for the configuration file")
	flag.Parse()
	urls := loadConfig(confPtr)
	for _, c := range urls.Checks {
		fmt.Println("Response time for:", c)
		measure(c)
	}
}

func main() {
	s := gocron.NewScheduler()
	s.Every(20).Seconds().Do(check)
	<-s.Start()
}
