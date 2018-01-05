package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/clsung/golami"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("Usage: ./golami <srv> <text>")
	}
	service := flag.Arg(0)
	text := flag.Arg(1)
	appKey := os.Getenv("OLAMI_APP_KEY")
	appSecret := os.Getenv("OLAMI_APP_SECRET")

	c, err := golami.New(appKey, appSecret)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.Post(context.Background(), service, text)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	/*
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	*/
	r := &golami.Result{}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&r); err != nil {
		log.Fatal(err)
	}
	switch service {
	case "seg":
		fmt.Println(r.Data.SEG)
	case "nli":
		//fmt.Printf("%v\n", r.Data.NLI)
		for _, n := range r.Data.NLI {
			fmt.Println(n.DescObj.Result)
		}
	}
}
