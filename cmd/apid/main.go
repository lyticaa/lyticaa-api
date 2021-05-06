package main

import (
	"os"
	"os/signal"

	"github.com/lyticaa/lyticaa-api/internal"
)

func main() {
	a := api.NewApi()
	a.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	a.Stop()
}
