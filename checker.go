package main

import "net/http"
import "time"
import "log"
import "fmt"
import "io/ioutil"
import "encoding/json"

type Checks struct {
	Urls []string
}

func loadConfig(path string) Checks {
	f, err := ioutil.ReadFile(path)

	var c Checks
	err = json.Unmarshal(f, &c)

	if err != nil {
		log.Fatal(err)
	}

	return c
}

func measure(address string) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := netClient.Get(address)

	if err != nil {
		panic(err)
	}

	start := time.Now()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Println("Load Time:", time.Since(start))
}

func main() {
	c := loadConfig("./config.json")
	for _, u := range c.Urls {
		fmt.Println("Response time for:", u)
		measure(u)
	}
}
