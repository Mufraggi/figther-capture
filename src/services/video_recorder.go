package services

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

const videoDuration = 6 * time.Minute // Durée de 6 minutes
func createFileName() string {
	timestamp := time.Now()
	dateStr := timestamp.Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s.mp4", dateStr)
}

func RecordVideo(webcam *gocv.VideoCapture, videoFile string) {
	window := gocv.NewWindow("Press 'Enter' to start recording")
	img := gocv.NewMat()
	defer window.Close()
	defer img.Close()

	// Attendre que l'utilisateur appuie sur Entrée pour démarrer l'enregistrement
	fmt.Println("Appuyez sur Entrée pour démarrer l'enregistrement vidéo de 6 minutes.")
	var input string
	fmt.Scanln(&input) // Attend l'entrée de l'utilisateur

	// Définir le writer pour enregistrer la vidéo1920x1080
	writer, _ := gocv.VideoWriterFile(videoFile, "h264", 24, 1920, 1080, true)
	defer writer.Close()

	// Capturer et enregistrer la vidéo pendant la durée spécifiée
	start := time.Now()
	fmt.Print(start)
	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			break
		}

		// Montrer l'image à l'écran
		window.IMShow(img)

		// Écrire dans le fichier vidéo
		writer.Write(img)

		// Sortir après 6 minutes
		if time.Since(start) > videoDuration {
			break
		}

		if window.WaitKey(1) == 27 { // Sortir avec Échap si besoin
			break
		}
	}
	fmt.Println("Enregistrement terminé")
}
