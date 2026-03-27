package main

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func Compress(file string) (string, error) {
	input, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer input.Close()

	outputFile := file + ".gz"
	output, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer output.Close()

	gzipWriter := gzip.NewWriter(output)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, input)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func CleanupOldBackups(dir string, max int) {
	files, _ := filepath.Glob(dir + "/*.gz")

	if len(files) <= max {
		return
	}

	sort.Strings(files)

	for i := 0; i < len(files)-max; i++ {
		os.Remove(files[i])
	}
}
