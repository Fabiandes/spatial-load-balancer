package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/fabiandes/spatial-load-balancer/api"
	"github.com/fabiandes/spatial-load-balancer/simulation"
	"github.com/fabiandes/spatial-load-balancer/simulation/entity"
	"go.uber.org/zap"
)

const (
	StartingEntityCount = 10
	SimulationRate      = 24
)

func main() {
	// Setup logging.
	cfg := zap.NewDevelopmentConfig()
	//cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	logger, _ := cfg.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	ctx := context.Background()

	// Set a random seed.
	rand.Seed(time.Now().UnixNano())

	// Configure and create simulation.
	opts := &simulation.Options{
		StartingEntityCount: StartingEntityCount,
		SimulationRate:      SimulationRate,
		Logger:              sugar,
	}

	s, err := simulation.New(opts)
	if err != nil {
		sugar.Errorf("Failed to create simulation: %v", err)
		return
	}

	api := api.New(&api.APIOptions{
		Logger: sugar,
	})
	go func() {
		if err := api.Listen(ctx, ":8080"); err != nil {
			sugar.Errorf("API failed to listen: %v", err)
		}
	}()
	go func() {
		ch := make(chan []*entity.Entity)
		s.Subscribe(ch)
		for es := range ch {
			if err := api.Broadcast(es); err != nil {
				sugar.Errorf("Failed to broadcast entities to clients: %v", err)
			}
		}
	}()

	if err := s.Run(ctx); err != nil {
		sugar.Errorf("An error occurred while running the simulation: %v", err)
		return
	}
}
