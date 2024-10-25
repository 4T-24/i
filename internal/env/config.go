package env

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	Domain   string `env:"DOMAIN" envDefault:"ctf.4ts.fr"`
	NodePort int    `env:"NODE_PORT" envDefault:"6368"`

	CTFd struct {
		URL   string `env:"URL"`
		Token string `env:"TOKEN"`
	} `envPrefix:"CTFD_"`

	Token       string `env:"TOKEN"`
	GlobalToken string `env:"GLOBAL_TOKEN"`

	SigningKey string `env:"SIGNING_KEY"`
}

func Get() *Config {
	return &cfg
}
