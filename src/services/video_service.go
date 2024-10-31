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
		}
		switch input {
		case "":
			filename := v.rec()
			if filename != nil {
				err := v.deleteVideo(*filename)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}
}

func (v *VideoService) createFileName() string {
	timestamp := time.Now()
	dateStr := timestamp.Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s.mp4", dateStr)
}

func (v *VideoService) rec() *string {
	file := v.createFileName()
	webcam, err := v.videoRecorder.NewWebCam()

	if err != nil {
		log.Fatal(err)
		return nil
	}
	writer, err := v.videoRecorder.NewWriter(file)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	v.videoRecorder.Rec(writer, webcam)
	err = v.c.Send(file)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &file
}
func (s *VideoService) deleteVideo(path string) error {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("Erreur lors de la suppression du fichier :", err)
		return err
	}
	return nil
}

func NewVideoService(
	videoRecorder IVideoRecorder,
	c IClientHttp) IVideoService {
	return &VideoService{videoRecorder, c}
}
