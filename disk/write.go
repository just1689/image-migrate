package disk

import (
	"fmt"
	"log"
	"os"
)

func NewWriter(path string, in chan string) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	go func() {
		defer f.Close()
		for i := range in {
			if _, err := f.WriteString(i + "\n"); err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}
