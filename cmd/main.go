package main

import (
	"context"
	"server/internal/config"
	"server/internal/run"

	_ "net/http/pprof"
)

func main() {

	config := config.MustLoadConfig(".")
	run.Run(context.Background(), config)
}
