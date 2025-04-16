package main

import (
	"net/url"
)

func main() {
	decodedKeyword := "%E4%BD%A0%E5%A5%BD"
	decodedKeyword, err := url.QueryUnescape(decodedKeyword)
	if err != nil {
		panic(err)
	}
	println(decodedKeyword)
}
