package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	AccessKey string `json:"ak"`
	SecretKey string `json:"sk"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
	BindAddr  string `json:"bind_addr"`
}

type Service struct {
	config *Config
	mac    *Mac
}

func New(cfg *Config) *Service {
	return &Service{config: cfg, mac: NewMac(cfg.AccessKey, cfg.SecretKey)}
}

func (s *Service) Run(addr string) error {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", s.handle)
	return http.ListenAndServe(addr, mux)
}

func (s *Service) handle(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		w.WriteHeader(405)
		w.Write([]byte("only get"))
		return
	}

	path := req.URL.Path
	u := fmt.Sprintf("http://%s%s", s.config.Domain, path)
	token := s.mac.Sign([]byte(u))
	url := fmt.Sprintf("%s?token=", u, token)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		return
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func Run(cfg *Config) {
	s := New(cfg)
	fmt.Println(s.Run(cfg.BindAddr))
}

func main() {
	confPath := flag.String("f", "", "config file")
	flag.Parse()
	if *confPath == "" {
		flag.PrintDefaults()
		return
	}
	data, err := ioutil.ReadFile(*confPath)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		return
	}
	config := &Config{}

	json.Unmarshal(data, config)
	Run(config)
}
