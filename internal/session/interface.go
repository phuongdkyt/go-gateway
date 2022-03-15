package session

import (
	"phuongnd/gateway/internal/domain"
	"time"
)

type SessionInterface interface {
	Create(sessionKey string) (domain.Session, error)
	Get(sessionId string) (string, error)
	Refresh(sessionId string) error
	//Development only
	SetTimeout(sessionId string, timeout time.Duration) error

	CreateSessionRedis(session domain.Session) error
	UpdateExpireRedis(session domain.Session) error
	FindKeyByIdRedis(id string) (string, error)
}
