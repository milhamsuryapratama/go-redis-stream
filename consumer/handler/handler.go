package handler

import (
	"github.com/milhamsuryapratama/go-redis-stream/packages/event"
)

func HandlerFactory() func(t event.Type) Handler {
	return func(t event.Type) Handler {
		switch t {
		case event.TrxType:
			return NewTransactionHandler()
		default:
			return NewTransactionHandler()
		}
	}
}

type Handler interface {
	Handle(event event.Event) error
}
