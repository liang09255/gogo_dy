package global

type config struct {
	Mysql struct {
		Dsn string `json:"dsn"`
	} `json:"mysql"`
	JwtKey string `json:"jwtKey"`
}

var Config = &config{}
