package docker

import (
	"fmt"
	"os/exec"
)

func Pull(tag string) (err error) {

	app := "docker"
	args := []string{
		"pull",
		tag,
	}

	cmd := exec.Command(app, args...)
	//var b []byte
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(string(b))
	return
}

func Tag(from, to string) (err error) {
	app := "docker"
	args := []string{
		"tag",
		from,
		to,
	}
	fmt.Println(app, args)
	cmd := exec.Command(app, args...)
	//var b []byte
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(string(b))
	return
}

func Push(tag string) (err error) {

	app := "docker"
	args := []string{
		"push",
		tag,
	}
	cmd := exec.Command(app, args...)
	//var b []byte
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(string(b))
	return
}
