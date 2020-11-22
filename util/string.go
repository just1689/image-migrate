package util

import "strings"

func SplitStringChan(in string) chan string {
	result := make(chan string)
	go func() {
		a := strings.Split(in, " ")
		for _, next := range a {
			result <- next
		}
		close(result)
	}()
	return result
}
