package helper

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	RepoName     string
	RepoOwner    string
	ClientID     string
	ClientSecret string
}

func RandBranchName() string {
	rand.Seed(time.Now().Unix())
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GetOptsFromEnv() (Config, error) {
	missingEnvMsg := "Set required env %s"

	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		return Config{}, fmt.Errorf(missingEnvMsg, "REPO_NAME")
	}

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		return Config{}, fmt.Errorf(missingEnvMsg, "REPO_OWNER")
	}

	client_id := os.Getenv("CLIENT_ID")
	if client_id == "" {
		return Config{}, fmt.Errorf(missingEnvMsg, "CLIENT_ID")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return Config{}, fmt.Errorf(missingEnvMsg, "SECRET")
	}

	return Config{
		RepoName:     repo,
		RepoOwner:    owner,
		ClientID:     client_id,
		ClientSecret: secret,
	}, nil
}
