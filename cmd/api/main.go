package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var sourceDir string
var destinationDir string
var copyFlag bool

func init() {
	flag.StringVar(&sourceDir, "source", "", "Directory to copy from")
	flag.StringVar(&destinationDir, "destination", "", "Directory to copy to")
	flag.BoolVar(&copyFlag, "c", false, "Flag to copy directory")
	flag.Parse()
}

func main() {
	if sourceDir == "" || destinationDir == "" {
		fmt.Println("Please provide a source and destination directory")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !copyFlag {
		fmt.Println("Copy flag not provided, exiting...")
		os.Exit(1)
	}

	parent := filepath.Base(sourceDir)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the directory is "Main" or "Solved", skip it
		if info.IsDir() && (info.Name() == "Main" || info.Name() == "Solved") {
			return filepath.SkipDir
		}

		// Generate the new path for file at destination directory
		relativePath := path[len(sourceDir):]
		destPath := filepath.Join(destinationDir, parent, relativePath)

		// Check if it is a directory
		// If it is, create the directory at destination
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// It is a file, let's copy it
		return copyFile(path, destPath)
	})

	if err != nil {
		fmt.Printf("Error copying folder: %v\n", err)
	}
}

func copyFile(sourcePath, destPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source) // check first var for number of bytes copied
	return err
}
