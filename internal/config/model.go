package config

type Config struct {
	Server   Server
	Database Database
	Mail     Email
	Redis    Redis
}

type Redis struct {
	Addr string
	Pass string
}

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Email struct {
	Host     string
	Port     string
	User     string
	Password string
}
