package drain

import "io"

//go:generate go run github.com/qazz-shyper/website/common/errors/errorgen

type Drainer interface {
	AcknowledgeReceive(size int)
	Drain(reader io.Reader) error
}
