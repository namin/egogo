package main

import (
	  "context"
    "fmt"
    "os"
	  "github.com/google/go-github/github"
)

func main() {
	client := github.NewClient(nil)

	me := "namin"

	fmt.Printf("followers for %v\n", me)
	page := 0
	total := 0
	perPage := 100
	for {
		opt := &github.ListOptions{PerPage: perPage, Page: page}
		users, resp, err := client.Users.ListFollowers(context.Background(), me, opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
			break
		} else {
			for _, user := range users {
				total += 1
				fmt.Printf("%v\n", *user.Login)
			}
			page = resp.NextPage
			if (page == 0) {
				break
			}
		}
	}
	fmt.Printf("total: %v\n", total)
}
