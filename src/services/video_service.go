package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type IVideoService interface {
	Run()
}
type VideoService struct {
	videoRecorder IVideoRecorder
	c             IClientHttp
}

func (v *VideoService) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Appuyez sur Entrée pour démarrer/arrêter l'enregistrement")
	fmt.Println("Appuyez sur 'q' pour quitter")
	for {
		char, err := reader.ReadString('\n')
		input := char[:len(char)-1]
		if err != nil {
			fmt.Println("Erreur de lecture:", err)
			continue
		}
		switch input {
		case "":
			v.rec()
		}

	}
}

func (v *VideoService) createFileName() string {
	timestamp := time.Now()
	dateStr := timestamp.Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s.mp4", dateStr)
}

func (v *VideoService) rec() {
	file := v.createFileName()
	webcam, err := v.videoRecorder.NewWebCam()

	if err != nil {
		log.Fatal(err)
		return
	}
	writer, err := v.videoRecorder.NewWriter(file)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer writer.Close()
	defer webcam.Close()

	v.videoRecorder.Rec(writer, webcam)
	err = v.c.Send(file)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewVideoService(
	videoRecorder IVideoRecorder,
	c IClientHttp) IVideoService {
	return &VideoService{videoRecorder, c}
}
