package main

import (
	"github.com/Mufraggi/figther-capture/src/services"
	"github.com/Netflix/go-env"
	"log"
	"time"
)

type Env struct {
	BackendUrl string `env:"BACKEND_URL" required:"true"`
}

func getEnvConfig[T any]() (*T, error) {
	var environment T
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}
func main() {
	env, err := getEnvConfig[Env]()
	if err != nil {
		log.Fatal(err)
	}
	client := services.NewClientHttp(env.BackendUrl)
	videoRecorder := services.NewVideoRecorder(10*time.Second, 0)
	videoService := services.NewVideoService(videoRecorder, client)
	videoService.Run()
}
