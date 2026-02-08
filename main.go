package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"omniversify-keystroke-visualizer/config"
	"omniversify-keystroke-visualizer/gui"
	"omniversify-keystroke-visualizer/input"
)

func init() {
	os.Stderr.WriteString("LOG: main.go init called\n")
	runtime.LockOSThread()
}

func main() {
	os.Stderr.WriteString("LOG: Program starting main...\n")
	cfg, err := config.Load("config.ini")
	if err != nil {
		fmt.Printf("Warning: could not load config.ini: %v\n", err)
		fmt.Println("Using default settings")
	}

	listener, err := input.NewListener()
	if err != nil {
		fmt.Printf("Error: failed to create input listener: %v\n", err)
		os.Exit(1)
	}

	events := make(chan input.Event, 100)

	if err := listener.Start(events); err != nil {
		fmt.Printf("Error: failed to start input listener: %v\n", err)
		os.Exit(1)
	}

	overlay, err := gui.NewOverlay(cfg)
	if err != nil {
		fmt.Printf("Error: failed to create overlay: %v\n", err)
		os.Exit(1)
	}

	// Visibility test event
	go func() {
		time.Sleep(1 * time.Second)
		events <- input.Event{Type: input.KeyEvent, Value: "Ready!"}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			listener.Stop()
			overlay.Close()
			return
		case ev := <-events:
			os.Stderr.WriteString(fmt.Sprintf("LOG: Main received event: %s\n", ev.Value))
			overlay.Show(ev.Value)
		}
	}
}
