package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func randBranchName() string {
	rand.Seed(time.Now().Unix())
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func handleError(err error, code int) {
	if err != nil {
		fmt.Println(err)
		os.Exit(code)
	}
}

func getOptsFromEnv() (config, error) {
	missingEnvMsg := "Set required env %s"

	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		return config{}, fmt.Errorf(missingEnvMsg, "REPO_NAME")
	}

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		return config{}, fmt.Errorf(missingEnvMsg, "REPO_OWNER")
	}

	client_id := os.Getenv("CLIENT_ID")
	if client_id == "" {
		return config{}, fmt.Errorf(missingEnvMsg, "CLIENT_ID")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return config{}, fmt.Errorf(missingEnvMsg, "SECRET")
	}

	return config{
		RepoName:     repo,
		RepoOwner:    owner,
		ClientID:     client_id,
		ClientSecret: secret,
	}, nil
}
