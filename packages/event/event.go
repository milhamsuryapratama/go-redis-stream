package event

import (
	"encoding"
	"time"
)

type Type string

const (
	TrxType Type = "TrxType"
)

type Event interface {
	GetID() string
	GetType() Type
	GetDateTime() time.Time
	SetID(id string)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
