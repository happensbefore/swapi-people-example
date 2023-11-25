package swclient

type Response[T any] struct {
	Data     []T
	NextPage int
}
