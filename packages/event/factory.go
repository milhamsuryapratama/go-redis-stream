package event

import "fmt"

func New(t Type) (Event, error) {
	b := &Base{
		Type: t,
	}

	switch t {

	case TrxType:
		return &TransactionEvent{
			Base: b,
		}, nil

	}

	return nil, fmt.Errorf("type %v not supported", t)
}
