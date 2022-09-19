package base

const (
	DOCKER            string = "docker"
	CONTAINER_RUN     string = "run -d %s"
	CONTAINER_START   string = "start %s"
	CONTAINER_STOP    string = "stop %s"
	CONTAINER_RESTART string = "restart %s"
	CONTAINER_PS      string = "ps"
	CONTAINER_PS_ALL  string = "ps -a"
	CONTAINER_RM      string = "rm %s"
	CONTAINER_RM_F    string = "rm -f %s"
	CONTAINER_COMMIT  string = "commit -a %s -m %s %s %s"

	IMAGE_ALL    string = "images"
	IMAGE_SEARCH string = "search %s --no-trunc"
	IMAGE_PULL   string = "pull %s"
	IMAGE_PUSH   string = "push %s"
	IMAGE_BUILD  string = "build -f %s -t %s ."
	IMAGE_REMOVE string = "rmi %s"
)
