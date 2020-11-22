package disk

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func ReadAllFiles(root string) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		err := filepath.Walk(root,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					//Ignore
					return nil
				}
				if !info.IsDir() {
					result <- path
				}
				//fmt.Println(path, info.Size())
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}()
	return result
}

func ReadFile(filename string) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
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
