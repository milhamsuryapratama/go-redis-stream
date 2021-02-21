package event

import (
	"encoding/json"
)

type TransactionEvent struct {
	*Base
	UserID int
	Amount int
}

func (t *TransactionEvent) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(t)
	return
}

func (t *TransactionEvent) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, t)
	return err
	//return msgpack.Unmarshal(data, t)
}
