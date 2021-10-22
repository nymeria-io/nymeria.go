package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"nymeria"
	"os"
)

var (
	auth  string
	help  bool
	purge bool
)

func getCacheDir() string {
	if dir := os.Getenv("NYMERIA_CACHE_DIR"); dir != "" {
		return dir
	}

	if dir, err := os.UserCacheDir(); err == nil {
		return fmt.Sprintf("%s/nymeria.io", dir)
	}

	return "/tmp/nymeria.io"
}

func purgeUserData() {
	os.RemoveAll(getCacheDir())
}

func cacheAuthKey(s string) {
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0750)
	if err := ioutil.WriteFile(fmt.Sprintf("%s/auth.key", cacheDir), []byte(s), 0600); err != nil {
		log.Println(err)
	}
}

func tryAuthFromCache() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/auth.key", getCacheDir()))

	if err != nil {
		return ""
	}

	return string(b)
}

func main() {
	flag.BoolVar(&help, "help", false, "Displays the tool's usage.")
	flag.BoolVar(&purge, "purge", false, "Purge all of the tool's cached data.")
	flag.StringVar(&auth, "auth", "", "Set's the tool's auth key. This will be be cached for future uses.")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if purge {
		purgeUserData()
		return
	}

	if len(auth) > 0 {
		cacheAuthKey(auth)

		if err := nymeria.SetAuth(auth); err != nil {
			log.Fatal(err)
		}
	} else {
		auth = tryAuthFromCache()
	}

	if len(auth) == 0 {
		fmt.Println("error: no auth key found")
		flag.Usage()
	}
}
