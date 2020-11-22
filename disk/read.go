package disk

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func ReadAllFiles(root string, recursive bool) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		if recursive {
			result <- root
			return
		}
		if err := filepath.Walk(root, fileHandler(result)); err != nil {
			log.Println(err)
		}
	}()
	return result
}

func fileHandler(result chan string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//Ignore
			return nil
		}
		if !info.IsDir() {
			result <- path
		}
		return nil
	}
}

func ReadFile(filename string) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		file, err := os.Open(filename)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}()
	return result
}
