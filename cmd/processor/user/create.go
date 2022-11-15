package user

import (
	"context"
	"math/rand"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	procerrors "pu/cmd/processor/errors"
	"pu/cmd/processor/user/domain"
	"pu/cmd/utils"
)

const saltCode = "!r#2@uM"

func (p *Processor) Create(ctx context.Context, u *domain.User, password string) (ru *domain.User, err error) {
	l := utils.GetLoggerFromContext(ctx)

	customSalt := generateCustomSalt()

	var passHash []byte
	if passHash, err = bcrypt.GenerateFromPassword([]byte(customSalt+password+saltCode), bcrypt.DefaultCost); err != nil {
		l.Error("can't generate password", zap.Error(err))
		return nil, procerrors.ErrGeneratePassword
	}

	u.PasswordHash = string(passHash)
	u.Salt = customSalt

	if ru, err = p.userRepo.Create(ctx, u); err != nil {
		l.Error("can't create user", zap.Error(err))
		return nil, procerrors.ErrDBRequest
	}

	if err = p.publisher.SendUserCreatedEvent(utils.GetTokenFromContext(ctx), ru); err != nil {
		l.Error("can't send event", zap.Error(err))
		return nil, procerrors.ErrPublishEvent
	}

	return ru, nil
}

func generateCustomSalt() string {
	return randSeq(10)
}

var letters = []rune("@#$!*&()%1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
