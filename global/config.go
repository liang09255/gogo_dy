package global

type config struct {
	Mysql struct {
		Dsn string `json:"dsn"`
	} `json:"mysql"`
	JwtKey string `json:"jwtKey"`
	AliOSS struct {
		Bucket          string `json:"bucket"`
		Endpoint        string `json:"endpoint"`
		AccessKeyId     string `json:"accessKeyId"`
		AccessKeySecret string `json:"accessKeySecret"`
	} `json:"aliOSS"`
	PasswordSalt string `json:"passwordSalt"`
}

var Config = &config{}
