package app

import (
	"context"
	"os"
	"time"

	"swapi/internal/infra/httpclient"
	"swapi/internal/infra/printer"
	"swapi/internal/infra/queue"
	"swapi/internal/infra/stringserializer"
	"swapi/internal/infra/swclient"
	"swapi/internal/infra/workerpool"
	"swapi/internal/services"
	"swapi/internal/services/flusher"
)

type Runnable interface {
	Run(ctx context.Context) error
}

const (
	swapiBaseEndpoint = "https://swapi.dev/api"
	swapiPersonMethod = "/people"
	swapiPaginator    = "/?page="
)

type App struct {
	services   []Runnable
	workerPool *workerpool.WorkerPool

	finisher chan struct{}
}

func New() *App {
	finisher := make(chan struct{})

	httpClient := httpclient.New(20 * time.Second)

	swapiClient := swclient.New(
		swclient.Config{
			BaseEndpoint: swapiBaseEndpoint,
			PersonMethod: swapiPersonMethod,
			Paginator:    swapiPaginator,
		},
		httpClient,
	)

	q := queue.NewQueue[string]()

	queueService := services.NewQueueService(q)

	personService := services.NewPersonService(finisher, swapiClient, queueService)

	printService := services.NewPrintService(stringserializer.New(), printer.New(os.Stdout))

	timeoutFlusher := flusher.NewTimeoutFlusher(q, printService)
	countFlusher := flusher.NewCountFlusher(q, printService)

	workerPool := workerpool.New()

	return &App{
		services: []Runnable{
			countFlusher,
			timeoutFlusher,
			printService,
			personService,
			queueService,
		},
		workerPool: workerPool,
		finisher:   finisher,
	}
}

func (a *App) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	a.workerPool.AddBackgroundJob(func() {
		<-a.finisher

		cancel()
	})

	for _, service := range a.services {
		s := service

		a.workerPool.AddBackgroundJob(func() {
			err := s.Run(ctx)
			if err != nil {
				panic(err)
			}
		})
	}

	a.workerPool.Run()
}
