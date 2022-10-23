package base

const (
	DOCKER string = "docker"
	// run -d --privileged=true --restart=always --name hello-world hello-world:latest
	CONTAINER_RUN string = "run"
	// 	exec -it abcp/base:0.1 /bin/sh /bin/passwd.sh pwd
	CONTAINER_EXEC string = "exec"
	// start hello-world
	CONTAINER_START string = "start"
	// stop hello-world
	CONTAINER_STOP string = "stop"
	// restart hello-world
	CONTAINER_RESTART string = "restart"
	CONTAINER_PS      string = "ps"
	CONTAINER_PS_ALL  string = "ps -a"
	// rm -f hello-world
	CONTAINER_RM string = "rm"
	// commit -a abcp -m desc "[base image]" hello-world my-hello-world:0.1
	CONTAINER_COMMIT string = "commit"
	// logs hello-world
	CONTAINER_LOG string = "logs"

	IMAGES       string = "images"
	IMAGE_SEARCH string = "search --no-trunc"
	IMAGE_PULL   string = "pull"
	// push my-hello-world:0.1
	IMAGE_PUSH string = "push"
	// build -f Dockerfile -t my-hello-world:0.2 .
	IMAGE_BUILD string = "build"
	// rmi my-hello-world:0.2
	IMAGE_REMOVE string = "rmi"
)
