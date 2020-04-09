package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codecommit"
)

func hasbranch(client *codecommit.CodeCommit, repositoryName string, branchName string) bool {
	found := false
	err := client.ListBranchesPages(&codecommit.ListBranchesInput{RepositoryName: &repositoryName},
		func(page *codecommit.ListBranchesOutput, lastPage bool) bool {
			for _, branch := range page.Branches {
				if *branch == branchName {
					found = true
				}
			}

			return lastPage
		})

	if err != nil {
		panic(err)
	}

	return found
}

func listRepositories(client *codecommit.CodeCommit, branchName string) {
	var repos []string
	err := client.ListRepositoriesPages(&codecommit.ListRepositoriesInput{},
		func(page *codecommit.ListRepositoriesOutput, lastPage bool) bool {
			for _, repo := range page.Repositories {
				repos = append(repos, *repo.RepositoryName)
			}

			return lastPage
		})

	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		if hasbranch(client, repo, branchName) {
			fmt.Println(repo)
		}
	}
}

func main() {
	var branchName string

	flag.StringVar(&branchName, "b", "", "Branch that we're looking for")
	flag.Parse()

	if "" == branchName {
		fmt.Println("Missing Branch flag")
		os.Exit(1)
	}

	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

	svc := codecommit.New(sess)

	listRepositories(svc, branchName)
}
