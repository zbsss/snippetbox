package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	imageTag := flag.String("image-tag", "", "Image tag")
	flag.Parse()

	if *imageTag == "" {
		log.Fatal("--image-tag is required")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = client.Close()
	}()

	contextDir := client.Host().Directory(".")

	// Unit tests
	unitTests := client.Container().
		From("golang:1.21-alpine").
		WithDirectory("/app", contextDir, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"deploy/", "infra/", "scripts/"}})

	out, err := unitTests.
		WithWorkdir("/app").
		WithExec([]string{"go", "test", "-v", "-short", "./..."}).
		Stderr(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	// Build and publish Docker image
	image := fmt.Sprintf("zbsss/snippetbox:%s", *imageTag)
	ref, err := contextDir.
		DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: "deploy/docker/snippetbox/Dockerfile",
		}).
		Publish(ctx, image)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Published image %q to :%s\n", image, ref)
}
