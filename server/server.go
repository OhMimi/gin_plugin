package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultPort         = 8080
	defaultReadTimeout  = 30
	defaultWriteTimeout = 30
)

type Server struct {
	name   string
	engin  *gin.Engine
	config *Config
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

func defaultConfig() *Config {
	return &Config{
		Port:         defaultPort,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}

// InstallFn install func type
type InstallFn func(r *gin.Engine)

func New(name string, config *Config) *Server {
	r := gin.New()

	return &Server{
		name:   name,
		engin:  r,
		config: config,
	}
}

func (s *Server) GetName() string {
	return s.name
}

func (s *Server) SetName(name string) {
	s.name = name
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engin
}

func (s *Server) SetEngine(r *gin.Engine) {
	s.engin = r
}

func (s *Server) GetConfig() *Config {
	return s.config
}

func (s *Server) SetConfig(c *Config) {
	s.config = c
}

// Install server install hook on gin engin
func (s *Server) Install(hooks ...InstallFn) {
	for _, fn := range hooks {
		fn(s.engin)
	}
}

func (s *Server) Run() {
	// set server serve parameters
	if s.config == nil {
		s.config = defaultConfig()
	}

	if s.config.Port == 0 {
		s.config.Port = defaultPort
	}

	if s.config.ReadTimeout == 0 {
		s.config.ReadTimeout = defaultReadTimeout
	}

	if s.config.WriteTimeout == 0 {
		s.config.WriteTimeout = defaultWriteTimeout
	}

	portStr := strconv.Itoa(s.config.Port)

	server := &http.Server{
		Addr:           ":" + portStr,
		Handler:        s.engin,
		ReadTimeout:    time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("server run on port: ", portStr)
	// start server
	if err := server.ListenAndServe(); err != nil {
		msg := fmt.Sprintf("server run failed Err: %v\n", err)
		panic(msg)
	}
}
