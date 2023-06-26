package system

import "go.uber.org/zap"

type Worker struct {
	logger *zap.SugaredLogger
}

func NewWorker(logger *zap.SugaredLogger) *Worker {
	w := &Worker{
		logger: logger,
	}

	return w
}

func (w *Worker) Update() error {
	//w.logger.Infoln("Worker update called")
	return nil
}
