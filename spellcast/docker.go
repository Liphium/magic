package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func DockerTest() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/hello-world", image.PullOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "hello-world",
	}, nil, nil, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatalln(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalln(err)
		}
	case st := <-statusCh:
		if st.Error != nil {
			log.Println(st.Error.Message)
		}
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatalln(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
