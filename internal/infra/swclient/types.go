package swclient

type swApiResult[T any] struct {
	Next    *string `json:"next"`
	Results []T     `json:"results"`
}
