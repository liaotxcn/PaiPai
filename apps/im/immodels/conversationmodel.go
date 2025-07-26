package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the mongo.
func NewConversationModel(url, db, collection string) ConversationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

func MustConversationModel(url, db string) ConversationModel {
	return NewConversationModel(url, db, "conversation")
}
