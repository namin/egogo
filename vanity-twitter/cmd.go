package main

import (
    "fmt"
    "net/url"
    "os"
    "github.com/ChimeraCoder/anaconda"
)

func getenv(id string) string {
	v := os.Getenv(id)
	if v == "" {
		panic(fmt.Sprintf("%v is not set", id))
	}
	return v
}

func main() {
	api_key := getenv("TWITTER_API_KEY")
	api_secret := getenv("TWITTER_API_SECRET")
	access_token := getenv("TWITTER_ACCESS_TOKEN")
	access_token_secret := getenv("TWITTER_ACCESS_TOKEN_SECRET")

	anaconda.SetConsumerKey(api_key)
	anaconda.SetConsumerSecret(api_secret)
	api := anaconda.NewTwitterApi(access_token, access_token_secret)
	api.DisableThrottling()
	api.ReturnRateLimitError(true)

	me := "nadamin"

	v := url.Values{}
	v.Set("count", "5000")

	fmt.Printf("followers for %v\n", me)
	total := 0
	c, err := api.GetFollowersIds(v)
	if err != nil {
		panic(err)
	} else {
		for _, id := range c.Ids {
			total += 1
			fmt.Printf("%v\n", id)
		}
	}
	fmt.Printf("total: %v\n", total)
}
