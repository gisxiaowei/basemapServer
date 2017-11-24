package config

type Config struct {
	Server   Server
	Services []Service
}

type Server struct {
	Port int64
}

type Service struct {
	Name string
	Path string
}
