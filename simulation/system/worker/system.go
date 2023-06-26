package worker

import (
	"fmt"
	"time"

	"github.com/fabiandes/spatial-load-balancer/simulation/entity"
	"go.uber.org/zap"
)

type Worker struct {
	logger     *zap.SugaredLogger
	currentJob Job
}

func NewWorker(logger *zap.SugaredLogger) *Worker {
	w := &Worker{
		logger: logger,
	}

	return w
}

// Update ensures that the Worker has a Job and that the next task is performed.
func (w *Worker) Update(dt time.Duration, e *entity.Entity) error {
	// If the worker doesn't currently have a job, it must be assigned one.
	if w.currentJob == nil {
		// TODO: Access the Job Manager to get a new job.
		w.currentJob = WanderingJob()
	}

	if err := w.currentJob.PerformNextTask(dt, e); err != nil {
		return fmt.Errorf("failed to perform next task: %v", err)
	}

	return nil
}

// // Job finds a new job for the worker.
// func (w *Worker) AssignNewJob() {
// 	// TODO: If there are no jobs available then set to wandering
// }
