package main

import (
	"log/slog"
	"os"

	"github.com/jparkkennaby/docker-kafka-go/system"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	s, err := system.Initialize()
	if err != nil {
		slog.Error("failed to initialse system", slog.Any("error", err))
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		slog.Error("pod stopped unexpectedly", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("pod stopped gracefully")
}
