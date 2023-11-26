package services

import (
	"context"
	"fmt"

	"swapi/internal/infra/lists"
	"swapi/internal/infra/swclient"
	"swapi/internal/models"
)

type SWClient interface {
	GetPersons(ctx context.Context, req *swclient.Request, resp *swclient.Response[models.Person]) error
}

type DataBuffer interface {
	Add(data []string)
}

type PersonService struct {
	dataLoadedEmitter chan<- struct{}
	swClient          SWClient
	dataBuffer        DataBuffer

	nextPage int
}

func NewPersonService(
	dataLoadedEmitter chan<- struct{},
	swClient SWClient,
	dataBuffer DataBuffer,
) *PersonService {
	return &PersonService{
		dataLoadedEmitter: dataLoadedEmitter,
		swClient:          swClient,
		dataBuffer:        dataBuffer,
		nextPage:          1,
	}
}

func (s *PersonService) Run(ctx context.Context) error {
	for s.hasMore() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		req := swclient.Request{
			NextPage: s.nextPage,
		}

		res := swclient.Response[models.Person]{}

		err := s.swClient.GetPersons(ctx, &req, &res)
		if err != nil {
			return fmt.Errorf("can't load persons: %w", err)
		}

		if len(res.Data) == 0 {
			continue
		}

		data := lists.Map(res.Data, func(item models.Person) string {
			return item.Name
		})

		s.dataBuffer.Add(data)

		s.nextPage = res.NextPage
	}

	close(s.dataLoadedEmitter)

	return nil
}

func (s *PersonService) hasMore() bool {
	return s.nextPage > 0
}
