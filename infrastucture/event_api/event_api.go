package event_api

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/zsmartex/pkg/v2/infrastucture/kafka"
)

type EventAPI struct {
	applicationName string
	producer        *kafka.Producer
	jwtPrivateKey   *rsa.PrivateKey
}

type EventAPIPayload struct {
	Record interface{} `json:"record"`
}

func New(producer *kafka.Producer, applicationName string, jwtPrivateKey string) (*EventAPI, error) {
	secret, err := base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(secret)
	if err != nil {
		return nil, err
	}

	return &EventAPI{
		applicationName: applicationName,
		producer:        producer,
		jwtPrivateKey:   privateKey,
	}, nil
}

func (e *EventAPI) generateJWT(event_payload EventAPIPayload) (string, error) {
	jwtPayload := jwt.MapClaims{
		"iat":   time.Now().Unix(),
		"jti":   strconv.FormatInt(time.Now().Unix(), 10),
		"iss":   e.applicationName,
		"exp":   time.Now().UTC().Add(time.Hour).Unix(),
		"event": event_payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtPayload)

	return token.SignedString(e.jwtPrivateKey)
}

func (e *EventAPI) Notify(context context.Context, event_name string, event_payload EventAPIPayload) error {
	eventType := strings.Split(event_name, ".")[0]
	topic := fmt.Sprintf("%s.events.%s", e.applicationName, eventType)
	jwtToken, err := e.generateJWT(event_payload)
	if err != nil {
		return err
	}

	e.producer.ProduceWithKey(context, topic, strings.Replace(event_name, fmt.Sprintf("%s.", eventType), "", 1), jwtToken)

	return nil
}
