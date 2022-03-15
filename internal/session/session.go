package session

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"phuongnd/gateway/internal/domain"
	"time"
)

var ProviderSessionSet = wire.NewSet(NewSessionRepository)
var ErrSessionKeyNotFound = errors.New("session key not found")

type SessionRepositoryImpl struct {
	redis redis.UniversalClient
}

func NewSessionRepository(redis redis.UniversalClient) SessionInterface {
	return &SessionRepositoryImpl{redis: redis}
}

func (s *SessionRepositoryImpl) Create(key string) (domain.Session, error) {
	sessionId := uuid.New().String()
	session := domain.Session{Id: sessionId, Key: key, Timeout: viper.GetDuration("SECRET_SESSION_TIMEOUT")}
	if err := s.CreateSessionRedis(session); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (s *SessionRepositoryImpl) Get(sessionId string) (string, error) {
	key, err := s.FindKeyByIdRedis(sessionId)
	if errors.Is(err, ErrSessionKeyNotFound) {
		return "", err
	} else if err != nil {
		return "", errors.Wrap(err, "get by session id failed")
	}
	return key, nil
}

func (s *SessionRepositoryImpl) Refresh(sessionId string) error {
	session := domain.Session{
		Id:      sessionId,
		Timeout: viper.GetDuration("SECRET_SESSION_TIMEOUT"),
	}
	return s.UpdateExpireRedis(session)
}

func (s *SessionRepositoryImpl) SetTimeout(sessionId string, timeout time.Duration) error {
	session := domain.Session{
		Id:      sessionId,
		Timeout: timeout,
	}
	return s.UpdateExpireRedis(session)
}

func (s *SessionRepositoryImpl) CreateSessionRedis(session domain.Session) error {
	return s.redis.Set(context.Background(), "session:"+session.Id, session.Key, session.Timeout).Err()
}

func (s *SessionRepositoryImpl) UpdateExpireRedis(session domain.Session) error {
	err := s.redis.Expire(context.Background(), "session:"+session.Id, session.Timeout).Err()
	return err
}

func (s *SessionRepositoryImpl) FindKeyByIdRedis(id string) (string, error) {
	secret, err := s.redis.Get(context.Background(), "session:"+id).Result()
	if err == redis.Nil {
		return "", ErrSessionKeyNotFound
	} else if err != nil {
		return "", errors.Wrap(err, "session key look-up failed")
	}
	return secret, nil
}
