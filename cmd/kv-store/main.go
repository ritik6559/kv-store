package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ritik6559/kv-store/internal/store/kv"
	"github.com/ritik6559/kv-store/internal/store/middleware"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	base, err := kv.NewKeyValueStore(100)
	if err != nil {
		log.Fatal(err)
	}
	metrics := middleware.NewMetricsMiddleware(base)
	s := middleware.NewLoggingMiddleware(metrics, logger.With("store", "kv"))

	cmds := []Command{
		{Op: "INCR", Key: "counter"},
		{Op: "INCR", Key: "counter"},
		{Op: "SET", Key: "name", Value: "ritik"},
		{Op: "INCR", Key: "name"},
		{Op: "RENAME", Key: "name", Value: "user"},
		{Op: "POP", Key: "user"},
		{Op: "GET", Key: "user"},
		{Op: "DELETE", Key: "counter"},
		{Op: "LEN"},
	}

	runCommands(s, cmds)

	fmt.Printf("metrics: %+v\n", metrics.Stats())
}
