package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Config struct {
	Port int `yaml:"port" env:"WS_PORT" env-default:"8081"`
}

type MessageType string

const (
	MessageTypeMessage MessageType = "Message"
	MessageTypeLeave   MessageType = "Leave"
	MessageTypeRead    MessageType = "Read"
)

type Message struct {
	ClientActorMessage
	Timestamp      int64     `json:"timestamp"`
	UserID         uuid.UUID `json:"userId"`
	ConversationID uuid.UUID `json:"conversationId"`
}

type ClientActorMessage struct {
	Type    MessageType `json:"messageType"`
	Message string      `json:"message"`
}

type ServerToActorMessages struct {
	Messages []ServerToActorMessage `json:"messages"`
}

type ServerToActorMessage struct {
	SenderID    uuid.UUID   `json:"senderId"`
	Message     string      `json:"message"`
	CreatedAt   int64       `json:"createdAt"`
	MessageType MessageType `json:"messageType"`
}

func (mes Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(mes)
}

func (mes Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &mes)
}
