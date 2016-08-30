package main

import "net/http"
import "time"
import "log"
import "fmt"
import "io/ioutil"
import "encoding/json"

type Urls struct {
	Checks []string
}

func loadConfig(path string) Urls {
	f, err := ioutil.ReadFile(path)

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
	urls := loadConfig("./config.json")
	for _, c := range urls.Checks {
		fmt.Println("Response time for:", c)
		measure(c)
	}
}
