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
			v.multipleCapture(3)
		}
	}
}

func (v *VideoService) multipleCapture(captureNb int) {
	timestamp := time.Now()
	dateStr := timestamp.Format("2006-01-02_15-04-05")
	for i := 0; i < captureNb; i++ {
		filename := v.createFileName(i, dateStr)
		v.rec(filename)
	}
}

func (v *VideoService) createFileName(run int, dateStr string) string {
	return fmt.Sprintf("%s-part_%d.mp4", dateStr, run)
}

func (v *VideoService) rec(file string) {
	fmt.Println("%s\n", file)
	filename, err := v.videoRecorder.Rec(file)
	if err != nil {
		log.Println("Error record:", err)
		return
	}
	go func() {
		timestamp1 := time.Now()
		dateStr1 := timestamp1.Format("2006-01-02_15-04-05")
		fmt.Println(filename, dateStr1, err)
		err := v.c.Send(*filename)
		timestamp := time.Now()
		dateStr := timestamp.Format("2006-01-02_15-04-05")
		fmt.Println(filename, dateStr, err)
		if err != nil {
			log.Println("Error sending file:", err)
			return
		}
		err = v.deleteVideo(file)
		if err != nil {
			log.Println("Error deleting file:", err)
		}
	}()
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
