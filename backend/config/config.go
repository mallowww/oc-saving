package config

type AppConfig struct {
	App struct {
		Name      string `mapstructure:"name"`
		Port      int    `mapstructure:"port"`
		Namespace string `mapstructure:"namespace"`
	} `mapstructure:"app"`
	Log struct {
		Environment string `mapstructure:"environment"`
		Level       string `mapstructure:"level"`
	} `mapstructure:"log"`
	Postgresql struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		GroupID string   `mapstructure:"group_id"`
		Topic   string   `mapstructure:"topic"`
		Auth    bool     `mapstructure:"auth"`
		Sasl    struct {
			Username  string `mapstructure:"username"`
			Password  string `mapstructure:"password"`
			Mechanism string `mapstructure:"mechanism"`
			Protocol  string `mapstructure:"protocol"`
		} `mapstructure:"sasl"`
	} `mapstructure:"kafka"`
}
