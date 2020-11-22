package docker

import "strings"

func IsDockerImage(in string) bool {
	return strings.Contains(in, "/") && strings.Contains(in, ":") && (strings.Index(in, ":") != len(in)-1)
}

func IsDockerImageSquishy(in string) bool {
	return strings.Contains(in, ":") && (strings.Index(in, ":") != len(in)-1)
}
