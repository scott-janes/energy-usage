package dao

type Config struct {
	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
	} `mapstructure:"database"`
	Kafka struct {
		Producer struct {
			Host  string `mapstructure:"host"`
			Port  int    `mapstructure:"port"`
			Topic string `mapstructure:"topic"`
		} `mapstructure:"producer"`
		Consumer struct {
			Host    string   `mapstructure:"host"`
			Port    int      `mapstructure:"port"`
			GroupID string   `mapstructure:"groupId"`
			Topics  []string `mapstructure:"topics"`
		} `mapstructure:"consumer"`
	} `mapstructure:"kafka"`
	ServiceName string `mapstructure:"serviceName"`
	Octopus     struct {
		APIKey       string `mapstructure:"apiKey"`
		MPAN         string `mapstructure:"mpan"`
		SerialNumber string `mapstructure:"serialNumber"`
	} `mapstructure:"octopus"`
}
