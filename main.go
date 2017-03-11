package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var number = flag.Int("n", 0, "Comic number.")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Make sure our random number generator is seeded.
	rand.Seed(time.Now().Unix())

	// Create and run our application.
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	status := app.GtkApp.Run(os.Args)
	if status != 0 {
		log.Printf("exit status: %v", status)
	}
}
