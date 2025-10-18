package main

import "github.com/nathaponb/pon2go/api"

func main() {
	err := api.NewApi()
	if err != nil {
		panic(err)
	}
}
