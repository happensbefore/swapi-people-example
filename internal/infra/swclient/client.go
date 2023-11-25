package swclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"swapi/internal/models"
)

type HTTPClient interface {
	Get(ctx context.Context, apiURL string) ([]byte, error)
}

type Client struct {
	cfg        Config
	httpClient HTTPClient
}

func New(cfg Config, httpClient HTTPClient) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (c *Client) GetPersons(ctx context.Context, req *Request, resp *Response[models.Person]) error {
	nextPage := req.NextPage
	if nextPage == 0 {
		nextPage = 1
	}

	basePath := fmt.Sprintf("%s%s%s", c.cfg.BaseEndpoint, c.cfg.PersonMethod, c.cfg.Paginator)
	apiUrl := fmt.Sprintf("%s%d", basePath, nextPage)

	data, err := c.httpClient.Get(ctx, apiUrl)
	if err != nil {
		return fmt.Errorf("can't do get persons request: %w", err)
	}

	swResult := swApiResult[models.Person]{}

	err = json.Unmarshal(data, &swResult)
	if err != nil {
		return fmt.Errorf("can't unmarshal get persons result: %w", err)
	}

	if swResult.Next != nil {
		next, ok := strings.CutPrefix(*swResult.Next, basePath)
		if ok {
			resp.NextPage, err = strconv.Atoi(next)
			if err != nil {
				return fmt.Errorf("can't unmarshal get persons result: %w", err)
			}
		}
	}

	resp.Data = swResult.Results

	return nil
}
