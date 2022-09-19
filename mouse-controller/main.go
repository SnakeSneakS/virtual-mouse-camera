package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/cmd"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatalf("2 arguments are required: addr port\ne.g.) go run main.go localhost 8080")
	}

	addr := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln(err)
	}

	cmd.Run(addr, port)
}
