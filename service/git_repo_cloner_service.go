package service

import (
	"errors"
	"log"
	"os/exec"
	"time"
)

func cloneRepo(repoUrl string, repoName string, targetPath string) error {
	folderExists, _ := folderExists(targetPath)
	folderIsEmpty, _ := isFolderEmpty(targetPath)

	if folderExists && !folderIsEmpty {
		log.Printf("INFO: %s has already been cloned. Skipping", repoName)
		return nil
	}

	log.Printf("INFO: Cloning %s", repoName)

	cmd := exec.Command("git", "clone", repoUrl, targetPath)

	err := cmd.Run()

	if err != nil {
		maxRetries := 5
		secondsToWait := 10
		numberOfRetries := 1
		for numberOfRetries <= maxRetries && err != nil {
			log.Printf("WARN: Cloning failed for: %s. Retrying...", repoName)
			time.Sleep(time.Duration(secondsToWait) * time.Second)
			err = cmd.Run()
			secondsToWait += 5
			numberOfRetries++
		}
		if numberOfRetries > maxRetries {
			log.Printf("WARN: Cloning of %s keeps failing. Skipping...", repoName)
			return errors.New("cloning failed")
		}
	} else {
		log.Printf("INFO: Cloning of '%s' succeeded", repoName)
	}
	return nil
}
