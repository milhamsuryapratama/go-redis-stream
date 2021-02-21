package handler

import (
	"fmt"
	"github.com/milhamsuryapratama/go-redis-stream/packages/event"
)

type TransactionHandler struct{}

func NewTransactionHandler() Handler {
	return &TransactionHandler{}
}

func (t *TransactionHandler) Handle(e event.Event) error {
	transactionEvent, ok := e.(*event.TransactionEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	fmt.Printf("processed event %+v UserID: %v Amount:%v \n", transactionEvent, transactionEvent.UserID, transactionEvent.Amount)

	return nil
}
