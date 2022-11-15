package user

import (
	"pu/cmd/database"
	"pu/cmd/publisher"
)

type (
	Processor struct {
		userRepo  *database.UserRepo
		publisher *publisher.Publisher
	}
)

func NewProcessor(
	userRepo *database.UserRepo,
	publisher *publisher.Publisher,
) *Processor {
	return &Processor{userRepo: userRepo, publisher: publisher}
}
