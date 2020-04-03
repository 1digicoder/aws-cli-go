
package main

import (
    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediastoredata"
	"github.com/jessevdk/go-flags"
	"fmt"
	"os"
	"strings"
)

var opts struct {
	Profile string `short:"p" long:"profile" description:"AWS Profile"`
	Region string `short:"r" long:"region" description:"AWS Region" default:"us-east-1"`
	Endpoint string `short:"e" long:"endpoint" description:"MediaStore Endpoint" required:"true"`
	Folder string `short:"f" long:"folder" description:"MediaStore folder" required:"true"`
	ForceFlag bool `long:"force" description:"Force operation without prompting."`
}

func ask_for_confirmation(elements int) bool {
	var s string

	fmt.Printf("Confirm deletion of %d(s) items. (y/N): ", elements)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options {
		SharedConfigState: session.SharedConfigEnable,
		Profile: opts.Profile,
	}))

	data_client := mediastoredata.New(sess, &aws.Config{
		Region: aws.String(opts.Region),
		Endpoint: aws.String(opts.Endpoint),
	})

	pageNum := 0
	var items []string

	err = data_client.ListItemsPages(
		&mediastoredata.ListItemsInput{Path: &opts.Folder},
		func(page *mediastoredata.ListItemsOutput, lastPage bool) bool {
			pageNum++
			for e := 0; e < len(page.Items); e++ {
				items = append(items, *page.Items[e].Name)
			}
			return lastPage == false
		})
	if err != nil {
		panic(err)
	}

	if len(items) == 0 {
		fmt.Println("No items to delete, exiting.")
		os.Exit(0)
	}

	if !opts.ForceFlag {
		ok_to_continue := ask_for_confirmation(len(items))
		if !ok_to_continue {
			fmt.Println("Cancelling operation!")
			os.Exit(1)
		}
	}

	var fullPath string
	for i := 0; i < len(items); i++ {
		fullPath = fmt.Sprintf("%s/%s", opts.Folder, items[i])
		fmt.Printf("Removing %s...\n",fullPath)
		_, err := data_client.DeleteObject(
			&mediastoredata.DeleteObjectInput{
				Path: aws.String(fullPath),
			})

		if err != nil {
			panic(err)
		}
	}
}