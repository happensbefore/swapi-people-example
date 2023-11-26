package swclient

import "swapi/internal/infra/lists"

type swApiResult[T any] struct {
	Next    *string       `json:"next"`
	Results lists.List[T] `json:"results"`
}
