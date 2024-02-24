package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

const (
	snippetboxImageRepository = "zbsss/snippetbox"
)

var platforms = []dagger.Platform{
	"linux/amd64",
	"linux/arm64",
}

func main() {
	imageTag := flag.String("image-tag", "", "Image tag")
	postsubmit := flag.Bool("postsubmit", false, "Indicates if this pipeline is part of postsubmit")
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

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		// Unit tests
		unitTests := client.Container(dagger.ContainerOpts{Platform: platform}).
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
		container := contextDir.DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: "deploy/docker/snippetbox/Dockerfile",
			Platform:   platform,
		})

		platformVariants = append(platformVariants, container)
	}

	var imageRef string
	if *postsubmit {
		imageRef = fmt.Sprintf("%s:latest", snippetboxImageRepository)
	} else {
		imageRef = fmt.Sprintf("%s:%s", snippetboxImageRepository, *imageTag)
	}

	imageDigest, err := client.Container().
		Publish(ctx, imageRef,
			dagger.ContainerPublishOpts{
				PlatformVariants: platformVariants,
			})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pushed multi-platform image w/ digest: ", imageDigest)
}
