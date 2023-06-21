package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/fabiandes/spatial-load-balancer/simulation"
	"go.uber.org/zap"
)

func main() {
	// Setup logging.
	cfg := zap.NewDevelopmentConfig()
	//cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	logger, _ := cfg.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Set a random seed.
	rand.Seed(time.Now().UnixNano())

	// Start profiling
	f, err := os.Create("myprogram.prof")
	if err != nil {
		sugar.Errorf("Failed to start profiling: %v", err)
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Configure and create simulation.
	opts := &simulation.Options{
		StartingEntityCount: 1000,
		Logger:              sugar,
	}

	s, err := simulation.New(opts)
	if err != nil {
		sugar.Errorf("Failed to create simulation: %v", err)
		return
	}

	// Cleanup
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		cancel()
	}()

	if err := s.Run(ctx); err != nil {
		sugar.Errorf("An error occurred while running the simulation: %v", err)
		return
	}
}
