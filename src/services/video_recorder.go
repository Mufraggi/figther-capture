package services

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

type IVideoRecorder interface {
	Rec(*gocv.VideoWriter, *gocv.VideoCapture)
	NewWriter(filename string) (*gocv.VideoWriter, error)
	NewWebCam() (*gocv.VideoCapture, error)
}
type VideoRecorder struct {
	duration   time.Duration
	cameraPort int
}

func (v *VideoRecorder) Rec(writer *gocv.VideoWriter, webcam *gocv.VideoCapture) {
	img := gocv.NewMat()
	defer img.Close()
	start := time.Now()
	fmt.Print("start at: %v", start)
	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			break
		}
		writer.Write(img)
		if time.Since(start) > v.duration {
			break
		}
	}
	defer writer.Close()
	defer webcam.Close()
	fmt.Print("end at: %v", time.Now())
}
func (v *VideoRecorder) NewWebCam() (*gocv.VideoCapture, error) {
	return gocv.OpenVideoCapture(v.cameraPort)
}

func (v *VideoRecorder) NewWriter(videoFile string) (*gocv.VideoWriter, error) {
	writer, err := gocv.VideoWriterFile(videoFile, "mp4v", 24, 1920, 1080, true)
	return writer, err
}

func NewVideoRecorder(duration time.Duration, cameraPort int) IVideoRecorder {
	return &VideoRecorder{duration: duration, cameraPort: cameraPort}
}
