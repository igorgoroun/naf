package ftnlib

import (
	"go/types"
	"time"
)

// FtnPacket - Anarchived fts-0001 packet
type FtnPacket struct {
	from     FtnAddress
	to       FtnAddress
	date     time.Time
	password string
	messages types.Array
}
