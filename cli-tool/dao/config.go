package dao

type Config struct {
	Database Database `mapstructure:"database"`
	Kafka    Kafka    `mapstructure:"kafka"`
}

type Database struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
}

type Kafka struct {
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`
  Topics []string `mapstructure:"TOPICS"`
}
