package publisher

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"pu/cmd/processor/user/domain"
	"pu/rmq"
)

// Publisher send events to broker
type Publisher struct {
	rc rmq.Client
}

func NewPublisher(client rmq.Client) *Publisher {
	return &Publisher{client}
}

type SendUserCreatedMessage struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Points    int64     `json:"points"`
}

func convertFromProcessor(user *domain.User) *SendUserCreatedMessage {
	return &SendUserCreatedMessage{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Points:    user.Points,
	}
}

func (p *Publisher) SendUserCreatedEvent(token string, user *domain.User) (err error) {
	msg := convertFromProcessor(user)
	var uj []byte
	if uj, err = json.Marshal(msg); err != nil {
		return errors.Wrap(err, "can't marshal event")
	}
	if err = p.rc.Publish(token, exchanger, routingKeyCreated, uj); err != nil {
		return errors.Wrap(err, "can't publish")
	}
	return nil
}
