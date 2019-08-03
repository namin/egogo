package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "github.com/ChimeraCoder/anaconda"
)

func getenv(id string) string {
	v := os.Getenv(id)
	if v == "" {
		panic(fmt.Sprintf("%v is not set", id))
	}
	return v
}

func get_map(map_file string)map[string]string {
	f, err := os.Open(map_file)
	if err != nil {
		panic(fmt.Sprintf("cannot read %v: %v", map_file, err))
	}
	defer f.Close()
	reader := csv.NewReader(f)
	m := make(map[string]string)
	first := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(fmt.Sprintf("error reading: %v", err))
		}
		if !first {
			m[record[0]] = record[1]
		}
		first = false
	}
	return m
}

func get_ids(id_file string)map[int]string {
	f, err := os.Open(id_file)
	if err != nil {
		panic(fmt.Sprintf("cannot read %v: %v", id_file, err))
	}
	defer f.Close()

	reader := csv.NewReader(f)
	m := make(map[int]string)
	total := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(fmt.Sprintf("error reading: %v", err))
		}
		m[total] = record[0]
		total += 1
	}
	return m
}

func write_ids(sn_file string, ids map[int]string, m map[string]string) {
	f, err := os.Create(sn_file)
	if err != nil {
		panic(fmt.Sprintf("cannot write %v: %v", sn_file, err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	n := len(ids)
	for i:=0; i < n; i++ {
		id := ids[i]
		if i==0 || i==n-1 {
			fmt.Fprintf(f, "%v\n", id)
		} else {
			sn := m[id]
			fmt.Fprintf(f, "%v\n", sn)
		}
	}
}

func map_keys(m map[string]string)[]string {
	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	return mk
}

func write_map(map_file string, m map[string]string) {
	f, err := os.Create(map_file)
	if err != nil {
		panic(fmt.Sprintf("cannot write %v: %v", map_file, err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Fprintf(f, "%v,%v\n", "id", "screen_name")
	mk := map_keys(m)
	for _, k := range mk {
		v := m[k]
		fmt.Fprintf(f, "%v,%v\n", k, v)
	}
}

func fetch_ids(api *anaconda.TwitterApi, a []int64, m map[string]string) {
	users, err := api.GetUsersLookupByIds(a, nil)
	if err != nil {
		panic(err)
	} else {
		for _, user := range users {
			m[user.IdStr] = user.ScreenName
		}
	}
}

func main() {
	api_key := getenv("TWITTER_API_KEY")
	api_secret := getenv("TWITTER_API_SECRET")
	access_token := getenv("TWITTER_ACCESS_TOKEN")
	access_token_secret := getenv("TWITTER_ACCESS_TOKEN_SECRET")

	anaconda.SetConsumerKey(api_key)
	anaconda.SetConsumerSecret(api_secret)
	api := anaconda.NewTwitterApi(access_token, access_token_secret)

	map_file := "twitter_map.csv"
	id_file := "twitter_ids.txt"
	sn_file := "twitter.txt"

  m := get_map(map_file)
	ids := get_ids(id_file)

	n := len(ids)
	j := 0
	a_n := 100
	var a [100]int64
	for i:=1; i<n-1; i++ {
		id := ids[i]
		_, ok := m[id]
		if !ok {
			intId, err := strconv.ParseInt(id, 0, 64)
			if err != nil {
				panic(err)
			}
			a[j] = intId
			j += 1
			if j==a_n {
				fetch_ids(api, a[0:100], m)
				j = 0
			}
		}
	}
	if j>0 {
		fetch_ids(api, a[0:j], m)
	}

	write_map(map_file, m)
	write_ids(sn_file, ids, m)
}
