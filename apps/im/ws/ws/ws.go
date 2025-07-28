package ws

import constants "PaiPai/pkg/constant"

type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string   `mapstructure:"sendId"`
		RecvId             string   `mapstructure:"recvId"`
		RecvIds            []string `mapstructure:"recvIds"`
		SendTime           int64    `mapstructure:"sendTime"`
		
		ReadRecords map[string]string     `json:"readRecords"`
		ContentType constants.ContentType `mapstructure:"contentType"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	// 已读消息
	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`
		RecvId             string   `mapstructure:"recvId"`
		ConversationId     string   `mapstructure:"conversationId"`
		MsgIds             []string `mapstructure:"msgIds"` // 消息id
	}
)
