package main

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/gorilla/securecookie"
)

func main() {
	var (
		n = flag.Int("n", 32, "length of the random key, not the generated string")
	)

	flag.Parse()

	bytes := securecookie.GenerateRandomKey(*n)
	if bytes == nil {
		panic("failed to generate random key")
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	fmt.Println("generate random string: ", encoded)
}
