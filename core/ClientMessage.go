package core

import "github.com/google/uuid"

type ClientMessage struct {
	ClientUuid uuid.UUID
	Payload    *[]byte
}
