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
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
)
