package main

import "net/http"
import "time"
import "log"
import "fmt"
import "flag"
import "io/ioutil"
import "encoding/json"

type Urls struct {
	Checks []string
}

func loadConfig(path string) Urls {
	confPtr := flag.String("config", "./config.json", "A fully qualified path for the configuration file")
	flag.Parse()

	f, err := ioutil.ReadFile(*confPtr)
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

func main() {
	urls := loadConfig("./config.json")
	for _, c := range urls.Checks {
		fmt.Println("Response time for:", c)
		measure(c)
	}
}
