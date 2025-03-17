package common

import "github.com/google/uuid"

type IdProvider struct{}

func (_ *IdProvider) Provide() uuid.UUID {
	return uuid.New()
}
