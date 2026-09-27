package main

import (
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
	s := middleware.NewLoggingMiddleware(base, logger.With("store", "kv"))

	cmds := []Command{
		{Op: "INCR", Key: "counter"},
		{Op: "INCR", Key: "counter"},
	}

	runCommand(s, cmds)
}
