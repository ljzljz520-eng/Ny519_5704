package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/sine3d/internal/httpapi"
	"example.com/sine3d/internal/preset"
)

func main() {
	address := flag.String("addr", ":8080", "HTTP listen address")
	fixture := flag.String("fixture", "", "optional YAML preset path")
	flag.Parse()

	defaultPreset, err := loadPreset(*fixture)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger := log.New(os.Stderr, "", 0)
	logger.Printf("sine3d listening on %s", *address)
	if err := http.ListenAndServe(*address, httpapi.New(defaultPreset)); err != nil {
		logger.Fatal(err)
	}
}

func loadPreset(path string) (preset.Preset, error) {
	if path == "" {
		return preset.LoadDefault()
	}
	return preset.LoadFile(path)
}
