package global

type config struct {
	Mysql struct {
		Dsn string `json:"dsn"`
	} `json:"mysql"`
}

var Config = &config{}
