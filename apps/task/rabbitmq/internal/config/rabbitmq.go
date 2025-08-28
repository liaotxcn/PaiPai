package config

type RabbitMQConfig struct {
	Host     string `json:",default=localhost"`
	Port     int    `json:",default=5672"`
	Username string `json:",default=guest"`
	Password string `json:",default=guest"`
	VHost    string `json:",default=/"`
}
