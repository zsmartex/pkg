package event_api_fx

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/infrastructure/kafka_fx"
)

var Module = fx.Module("event_api_fx", fx.Provide(New))

type EventAPI struct {
	producer        *kafka_fx.Producer
	jwtPrivateKey   *rsa.PrivateKey
	applicationName string
}

type EventAPIPayload struct {
	Record interface{} `json:"record"`
}

type eventAPIParams struct {
	fx.In

	Producer        *kafka_fx.Producer
	ApplicationName string `name:"application_name"`
	JWTPrivateKey   string `name:"event_api_jwt_private_key"`
}

func New(params eventAPIParams) (*EventAPI, error) {
	secret, err := base64.StdEncoding.DecodeString(params.JWTPrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(secret)
	if err != nil {
		return nil, err
	}

	return &EventAPI{
		producer:        params.Producer,
		jwtPrivateKey:   privateKey,
		applicationName: params.ApplicationName,
	}, nil
}

func (e *EventAPI) generateJWT(event_payload EventAPIPayload) (string, error) {
	jwtPayload := jwt.MapClaims{
		"iat":   time.Now().Unix(),
		"jti":   time.Now().Unix(),
		"iss":   e.applicationName,
		"exp":   time.Now().Add(time.Hour).Unix(),
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

	key := strings.Replace(event_name, fmt.Sprintf("%s.", eventType), "", 1)

	return e.producer.ProduceWithKey(context, topic, []byte(key), jwtToken)
}
