package main

import (
	"flag"
	"log"
	"net/http"

	"maskhub/internal/console"
	"maskhub/internal/schedule"
	"maskhub/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	data := flag.String("data", "./data", "data directory")
	flag.Parse()

	state := store.NewState(*data)
	handler := console.NewAPI(state)

	runner := schedule.NewRunner(state, schedule.DefaultOptions())
	runner.Start()
	defer runner.Stop()

	log.Printf("MaskHub listening on %s with data dir %s", *addr, *data)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}
