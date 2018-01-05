package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/clsung/golami"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("Usage: ./golami <srv> <text|file>")
	}
	service := flag.Arg(0)
	text := flag.Arg(1)
	appKey := os.Getenv("OLAMI_APP_KEY")
	appSecret := os.Getenv("OLAMI_APP_SECRET")

	c, err := golami.New(appKey, appSecret)
	if err != nil {
		log.Fatal(err)
	}
	var r *golami.Result
	if service == "asr" {
		r, err = c.PostASR(context.Background(), text)
	} else {
		r, err = c.PostText(context.Background(), service, text)
	}

	if err != nil {
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
	case "asr":
		fmt.Println(r.Data.ASR.Result)
		fmt.Println(r.Data.SEG)
		for _, n := range r.Data.NLI {
			fmt.Println(n.DescObj.Result)
		}
	}
}
