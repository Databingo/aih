package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
        var ghp_token string
	fmt.Println("Please input ghp_token")
	fmt.Scanf("%s", & ghp_token)
	fmt.Println("ghp_token:", ghp_token)

	ctx := context.Background()

	// Authenticate with GitHub using a personal access token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghp_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Set the owner and repository for the release
	owner := "Databingo"
	repo := "aih"

	// Create a new release
	release, _, err := client.Repositories.CreateRelease(ctx, owner, repo, &github.RepositoryRelease{
		TagName:         github.String("v0.1.0"),
		TargetCommitish: github.String("master"),
		Name:            github.String("Release v0.1.0"),
		Body:            github.String("Welcome to Aih, an open plan based on the idea of \"Co-relation's enhancement of AI and human beings\""),
		Draft:           github.Bool(false),
		Prerelease:      github.Bool(false),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload binary zip files to the release
	files := []string{ "aih_amd64_exe.zip", "aih_arm64_exe.zip", "aih_arm_exe.zip", "aih_linux_amd64.zip", "aih_linux_arm.zip", "aih_linux_arm64.zip", "aih_linux_x86.zip", "aih_mac_amd64.zip", "aih_mac_arm64.zip", "aih_x86_exe.zip"}
	for _, file := range files {
		// Read the file contents
		_, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		// Create the file on GitHub
		fileContent, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fileContent.Close()

		_, _, err = client.Repositories.UploadReleaseAsset(ctx, owner, repo, release.GetID(), &github.UploadOptions{
			Name: file,
		}, fileContent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Uploaded %s to release %d\n", file, release.GetID())
	}
}
