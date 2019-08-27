package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/dertseha/zeitraffer/internal/mswin"
)

func setupSignalHandler(quit chan<- struct{}) {
	signalChannel := make(chan os.Signal, 2)
	signals := []os.Signal{os.Interrupt, os.Kill}
	signal.Notify(signalChannel, signals...)
	go func() {
		<-signalChannel
		signal.Reset(signals...)
		close(quit)
	}()
}

func main() {
	runtime.LockOSThread()

	grabber, err := mswin.NewGrabber()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create grabber: %v\n", err)
		os.Exit(-1)
	}
	defer grabber.Dispose()

	quit := make(chan struct{})
	setupSignalHandler(quit)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := false
	counter := 0
	for !done {
		select {
		case <-quit:
			done = true
		case <-ticker.C:
			img := grabber.Grab()
			save(img, counter)
			counter++
		}
	}
}

func save(image image.Image, id int) {
	filename := fmt.Sprintf("%05d.png", id)
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	_ = png.Encode(file, image)
}
