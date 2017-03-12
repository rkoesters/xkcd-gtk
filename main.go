package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Make sure our random number generator is seeded.
	rand.Seed(time.Now().Unix())

	// Create and run our application.
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(app.GtkApp.Run(os.Args))
}
