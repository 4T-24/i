package env

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	Domain   string `env:"DOMAIN" envDefault:"ctf.4ts.fr"`
	NodePort int    `env:"NODE_PORT" envDefault:"6368"`
}

func Get() *Config {
	return &cfg
}
