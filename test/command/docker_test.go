package command_test

import (
	"CloudPlatform/base"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestGetImages(t *testing.T) {
	out, err := executor(base.DOCKER, base.IMAGE_ALL)
	if err != nil {
		t.Error(err)
	}
	str := string(out)
	fmt.Println(str)
	/*
	   REPOSITORY               TAG       IMAGE ID       CREATED      SIZE
	   achobeta/abcp_database   0.1       2952644bce9e   3 days ago   117MB
	*/
	ss := strings.Split(str, "\n")
	fmt.Println(ss[1])
	/*
	   achobeta/abcp_database   0.1       2952644bce9e   3 days ago   117MB
	*/
}

func TestSearchImages(t *testing.T) {
	cmd := fmt.Sprintf(base.IMAGE_SEARCH, "hello-world")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
	   NAME                                       DESCRIPTION                                     STARS     OFFICIAL   AUTOMATED
	   hello-world                                Hello World! (an example of minimal Dockeriz…   1848      [OK]
	   kitematic/hello-world-nginx                A light-weight nginx container that demonstr…   152
	   tutum/hello-world                          Image to test docker deployments. Has Apache…   89                   [OK]
	   dockercloud/hello-world                    Hello World!                                    19                   [OK]
	*/
}

func TestBuildImageByDockerfile(t *testing.T) {
	cmd := fmt.Sprintf(base.IMAGE_BUILD, "./TestDockerfile", "my-image:test")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
}

func TestRemoveImage(t *testing.T) {
	cmd := fmt.Sprintf(base.IMAGE_REMOVE, "my-image:test")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(out)
	/*
		Untagged: my-image:test
		Deleted: sha256:ed38e2b0124e1e76f2f19fd64272f87c8be87abf29c22964331e3eef30a750be
	*/
}

func TestPullImage(t *testing.T) {
	cmd := fmt.Sprintf(base.IMAGE_PULL, "hello-world:latest")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
}

// 案例会不通过
func TestPushImage(t *testing.T) {
	cmd := fmt.Sprintf(base.IMAGE_PUSH, "hello-world:latest")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		fmt.Println(string(out))
		t.Error(err)
	}
	fmt.Println(string(out))
}

func TestCreateContainer(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_RUN, "--name hello-world hello-world:latest")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		fmt.Println(string(out))
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
		98e6a4dd3e3b8912f156773f33adea70f6238a1ca3ee9bd18148d72ff9956c18
	*/
}

// 一般查容器是通过数据库，而不会通过 docker ps -a
func TestGetContainers(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_PS_ALL)
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
		CONTAINER ID   IMAGE                        COMMAND                  CREATED         STATUS                     PORTS                    NAMES
		98e6a4dd3e3b   hello-world:latest           "/hello"                 5 minutes ago   Exited (0) 5 minutes ago                            hello-world
		c20c91c90d34   hello-world                  "/hello"                 6 minutes ago   Exited (0) 6 minutes ago                            friendly_villani
		f0f941a603f1   hello-world                  "/hello"                 7 minutes ago   Exited (0) 7 minutes ago                            dazzling_borg
		a6d378d81c01   achobeta/abcp_database:0.1   "docker-entrypoint.s…"   7 hours ago     Up 7 hours                 0.0.0.0:6379->6379/tcp   abcp_database
	*/
}

func TestContainerStart(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_START, "hello-world")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
		hello-world
	*/
}

func TestContainerStop(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_STOP, "hello-world")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
		hello-world
	*/
}

func TestContainerRestart(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_RESTART, "hello-world")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
		hello-world
	*/
}

func TestMakeImageByContainer(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_COMMIT, "慢慢", "测试", "hello-world", "my-test:0.1")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(out)
	/*
		sha256:d7a966a74f16d7fa5f38078500e2b2cf072b83b9943aca23072ec2b212b275ea
	*/
}

func TestRemoveContainer(t *testing.T) {
	cmd := fmt.Sprintf(base.CONTAINER_RM, "hello-world")
	out, err := executor(base.DOCKER, cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
	/*
	   hello-world
	*/
}

// 在容器中传文件（待完善）
func TestUpload(t *testing.T) {

}

func executor(name, arg string) (string, error) {
	out, err := exec.Command(name, strings.Split(arg, " ")...).Output()
	return string(out), err
}
