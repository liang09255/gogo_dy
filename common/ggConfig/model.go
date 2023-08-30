package ggConfig

type config struct {
	Mysql struct {
		Username string
		Password string
		Host     string
		Port     int
		Db       string
	}
	JwtKey string
	AliOSS struct {
		Bucket          string
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
	}
	PasswordSalt string
	Redis        struct {
		Host     string
		Port     int
		Password string
		Db       int
	}
	Etcd struct {
		Addrs []string
	}
	UserServer  Server
	VideoServer Server
	ChatServer  Server
}

type Server struct {
	Name string
	Addr string
	Http string
	Port string
}
