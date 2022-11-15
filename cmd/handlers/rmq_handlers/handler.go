package rmqhandlers

import (
	rmqwrapper "github.com/naf-m62/rabbitmq_wrapper"

	"pu/cmd/processor/user"
)

type Handler struct {
	userProcessor *user.Processor
}

func NewHandler(
	userProcessor *user.Processor,
) *Handler {
	return &Handler{
		userProcessor: userProcessor,
	}
}

func RegisterHandlerList(h *Handler) []rmqwrapper.ConsumeItem {
	return []rmqwrapper.ConsumeItem{
		{
			ServiceName: "user",
			Exchange:    "order",
			RoutingKey:  "addPoints",
			Handler:     h.AddPoints,
		},
	}
}
