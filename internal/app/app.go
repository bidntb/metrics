package app

import (
	"log"
	"time"

	"bidntb/metrics/internal/middleware/handler"
	"bidntb/metrics/internal/nconfig"
	"bidntb/metrics/internal/router"
	"bidntb/metrics/internal/service/metrics"
	"bidntb/metrics/internal/storage"
)

func Run() {
	cfg := nconfig.ParseConfig()

	storageInstance := storage.NewMemStorage()
	metricsSvc := metrics.NewService(storageInstance)
	h := handler.NewHandler(metricsSvc)

	if cfg.Restore {
		if err := metricsSvc.LoadFrom(cfg.FilePath); err != nil {
			log.Printf("failed to restore metrics from %s: %v", cfg.FilePath, err)
		} else {
			log.Printf("metrics restored from %s", cfg.FilePath)
		}
	}

	var ticker *time.Ticker
	if cfg.StoreInterval > 0 {
		ticker = time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		go func() {
			for range ticker.C {
				if err := metricsSvc.SaveTo(cfg.FilePath); err != nil {
					log.Printf("periodic save error: %v", err)
				}
			}
		}()
	} else {
		log.Printf("STORE_INTERVAL=0, synchronous persistence mode (save on update if needed)")
	}

	r := router.SetupRouter(h)

	if err := r.Run(cfg.ServerAddr); err != nil {
		panic(err)
	}

	if ticker != nil {
		ticker.Stop()
	}
	if err := metricsSvc.SaveTo(cfg.FilePath); err != nil {
		log.Printf("final save error: %v", err)
	}
}
