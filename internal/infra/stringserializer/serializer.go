package stringserializer

import "fmt"

type Serializer struct{}

func New() *Serializer {
	return &Serializer{}
}

func (s *Serializer) Serialize(data string) []byte {
	return []byte(fmt.Sprintf("%s\n", data))
}
