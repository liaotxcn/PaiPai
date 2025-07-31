package configserver

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrNotSetConfig = errors.New("none config info...")

type OnChange func([]byte) error

type ConfigServer interface {
	Build() error
	SetOnChange(OnChange)
	FromJsonBytes() ([]byte, error)
}

type configSever struct {
	ConfigServer
	configFile string
}

func NewConfigServer(configFile string, s ConfigServer) *configSever {
	return &configSever{
		ConfigServer: s,
		configFile:   configFile,
	}
}

func (s *configSever) MustLoad(v any, onChange OnChange) error {
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}
	if s.ConfigServer == nil {
		// 使用go-zero默认方式
		conf.MustLoad(s.configFile, v)
		return nil
	}
	// 判断是否配置动态加载
	if onChange != nil {
		s.SetOnChange(onChange)
	}
	// 构建配置服务
	if err := s.ConfigServer.Build(); err != nil {
		return err
	}

	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(data, v)
}

func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
