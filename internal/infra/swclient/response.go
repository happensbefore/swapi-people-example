package swclient

import "swapi/internal/infra/lists"

type Response[T any] struct {
	Data     lists.List[T]
	NextPage int
}
