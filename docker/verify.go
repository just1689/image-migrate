package docker

import "strings"

func IsDockerImage(in string) bool {
	return strings.Contains(in, "/") && (strings.LastIndex(in, "/") < (strings.Index(in, ":"))) && IsDockerImageSquishy(in)
}

func IsDockerImageSquishy(in string) bool {
	return strings.Contains(in, ":") && (strings.Index(in, ":") != len(in)-1)
}
