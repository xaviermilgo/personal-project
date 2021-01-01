package main

import (
	crand "crypto/rand"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const startDate = "2021-01-01T00:00:00.00Z"

func initRepo() (err error) {
	log.Info("initRepo()")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("working dir: %s", workingDir)

	repo, err := git.PlainInit(workingDir, false)
	if err != nil {
		return
	}

	log.Info("Git init success")

	workTree, err := repo.Worktree()
	if err != nil {
		return
	}

	initialFiles := []string{"main.go", ".gitignore", "file", "go.mod", "go.sum"}

	for _, file := range initialFiles {
		_, err = workTree.Add(file)
		if err != nil {
			return
		}
	}

	start, err := time.Parse(time.RFC3339, startDate)

	if err != nil {
		return
	}

	commit, err := workTree.Commit("init", &git.CommitOptions{
		Author: &object.Signature{
			Name:  getVarFromEnv("GIT_UNAME"),
			Email: getVarFromEnv("GIT_EMAIL"),
			When:  start,
		},
	})

	_, err = repo.CommitObject(commit)
	if err != nil {
		return
	}

	return nil

}

func getVarFromEnv(envName string) (env string) {
	env, exists := os.LookupEnv(envName)
	if exists == false {
		err := fmt.Errorf("Env variable %s is not set", envName)
		panic(err)
	}
	return
}

func main() {

	log.Info("stimmoc_rekaf starting")

	err := godotenv.Load()

	if err != nil {
		log.Error(err)
		return

	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Error(err)
		return
	}

	repo, err := git.PlainOpen(workingDir)

	if err == git.ErrRepositoryNotExists {
		log.Warn("Repo not initialized, creating...")
		err = initRepo()
		if err == nil {
			repo, err = git.PlainOpen(workingDir)
		}
	}

	if err != nil {
		log.Error(err)
		return
	}

	ref, err := repo.Head()
	if err != nil {
		log.Error(err)
		return
	}

	lc, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Error(err)
		return
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return
	}

	startT := lc.Author.When
	stopT := time.Now()

	counter := startT

	log.Info("creating commits")

	rand.Seed(time.Now().UnixNano())

	minM := 60
	maxM := 720

	for {

		if stopT.Truncate(24 * time.Hour).Before((counter.Truncate(24 * time.Hour))) {
			log.Info("Commits up to date")
			return
		}

		r := rand.Intn(maxM-minM+1) + minM

		counter = counter.Add(time.Minute * time.Duration(r))

		if counter.Weekday() == time.Sunday {
			if rand.Intn(100) >= 10 {
				continue
			}
		}

		randomBytes := make([]byte, 16)
		_, err = crand.Read(randomBytes)
		if err != nil {
			log.Error(err)
			return
		}

		err = ioutil.WriteFile("file", randomBytes, 0644)
		if err != nil {
			log.Error(err)
			return
		}

		_, err = workTree.Add("file")
		if err != nil {
			log.Error(err)
			return
		}

		msg := fmt.Sprintf("%x", randomBytes)
		commit, err := workTree.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  getVarFromEnv("GIT_UNAME"),
				Email: getVarFromEnv("GIT_EMAIL"),
				When:  counter,
			},
		})

		_, err = repo.CommitObject(commit)
		if err != nil {
			log.Error(err)
			return
		}

	}

}
