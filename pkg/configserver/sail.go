package configserver

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/HYY-yu/sail-client"
)

// 修改ETCDEndpoints类型为字符串数组
type Config struct {
	ETCDEndpoints  []string `toml:"etcd_endpoints"`
	ProjectKey     string   `toml:"project_key"`
	Namespace      string   `toml:"namespace"`
	Configs        string   `toml:"configs"`
	ConfigFilePath string   `toml:"config_file_path"`
	LogLevel       string   `toml:"log_level"`
}

// 辅助函数：将字符串数组转换为逗号分隔的字符串
func endpointsToString(endpoints []string) string {
	return strings.Join(endpoints, ",")
}

// 辅助函数：将逗号分隔的字符串转换为字符串数组
func stringToEndpoints(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// 修改NewSail函数，适应ETCDEndpoints类型变更
func NewSail(cfg *Config) *Sail {
	s := sail.New(&sail.MetaConfig{
		ETCDEndpoints:  endpointsToString(cfg.ETCDEndpoints),
		ProjectKey:     cfg.ProjectKey,
		Namespace:      cfg.Namespace,
		Configs:        cfg.Configs,
		ConfigFilePath: cfg.ConfigFilePath,
		LogLevel:       cfg.LogLevel,
	})
	return &Sail{Sail: s, c: cfg}
}

type Sail struct {
	*sail.Sail
	sail.OnConfigChange
	c *Config
}

func (s *Sail) Build() error {
	var opts []sail.Option
	if s.OnConfigChange != nil {
		opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}
	s.Sail = sail.New(&sail.MetaConfig{
		ETCDEndpoints:  endpointsToString(s.c.ETCDEndpoints),
		ProjectKey:     s.c.ProjectKey,
		Namespace:      s.c.Namespace,
		Configs:        s.c.Configs,
		ConfigFilePath: s.c.ConfigFilePath,
		LogLevel:       s.c.LogLevel,
	}, opts...)
	return s.Err()
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}
	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) SetOnChange(f OnChange) {
	s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
		data, err := s.fromJsonBytes(sail)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = f(data); err != nil {
			fmt.Println("OnChange err:", err)
		}
	}
}
func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}
