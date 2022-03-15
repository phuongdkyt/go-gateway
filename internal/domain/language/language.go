package language

import (
	"context"
)

type Locale string

const (
	Vietnamese Locale = "vi"
	English    Locale = "en"
)

type UserMessage struct {
	Code    string
	Message map[Locale]string
}

type Repository interface {
	FetchAll(ctx context.Context) ([]UserMessage, error)
	Upsert(ctx context.Context, messages []UserMessage) error
}

type Service interface {
	//Load load all user messages into memory
	Load(ctx context.Context) error
	CreateOrUpdateMessageList(ctx context.Context, messages []UserMessage) error
	//Translate translate message code into end-user message. If code not match, this return the code itself.
	Translate(ctx context.Context, code string, loc Locale) string
}
