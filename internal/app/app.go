package app

import (
	"bidntb/metrics/internal/middleware/handler"
	"bidntb/metrics/internal/nconfig"
	"bidntb/metrics/internal/router"
	"bidntb/metrics/internal/service/metrics"
	"bidntb/metrics/internal/storage"
)

func Run() {
	serverAddress := nconfig.GetServerAddress()
	storageInstance := storage.NewMemStorage()
	metricsSvc := metrics.NewService(storageInstance)
	h := handler.NewHandler(metricsSvc)

	r := router.SetupRouter(h)

	if err := r.Run(serverAddress); err != nil {
		panic(err)
	}
}
