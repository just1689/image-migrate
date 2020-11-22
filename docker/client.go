package docker

import (
	"os/exec"
)

func Pull(tag string) (err error) {
	app := "docker"
	args := []string{
		"pull",
		tag,
	}
	return execute(app, args)
}

func Tag(from, to string) (err error) {
	app := "docker"
	args := []string{
		"tag",
		from,
		to,
	}
	return execute(app, args)
}

func Push(tag string) (err error) {

	app := "docker"
	args := []string{
		"push",
		tag,
	}
	return execute(app, args)
}

func execute(app string, args []string) (err error) {
	cmd := exec.Command(app, args...)
	_, err = cmd.Output()
	return
}
