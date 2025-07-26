package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversations struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	UserId           string                   `bson:"userId"`
	ConversationList map[string]*Conversation `bson:"conversationList"`

	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
