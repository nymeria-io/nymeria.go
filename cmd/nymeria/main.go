package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"git.nymeria.io/nymeria.go"
)

var (
	auth                string
	bulkenrich          string
	checkAuthentication bool
	enrich              string
	help                bool
	purge               bool
	verify              string
)

func prettyPrint(i interface{}) {
	if bs, err := json.MarshalIndent(i, "", "  "); err == nil {
		fmt.Println(string(bs))
	}
}

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
	flag.BoolVar(&checkAuthentication, "check-auth", false, "If set, will test the supplied or cached api key and determine if it's valid or not.")
	flag.StringVar(&auth, "auth", "", "Set's the tool's auth key. This will be be cached for future uses.")
	flag.StringVar(&verify, "verify", "", "If an email is specified, will try to discover the deliverability of the email using Nymeria's API.")
	flag.StringVar(&enrich, "enrich", "", "A JSON encoded set of enrichment params (ex: '[{'url': 'github.com/nymeriaio'}]')")

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
	} else {
		auth = tryAuthFromCache()
	}

	if len(auth) == 0 {
		fmt.Println("error: no auth key found")
		return
	}

	/* we have an api key, we can begin doing real work */

	if err := nymeria.SetAuth(auth); err != nil {
		log.Fatal(err)
	}

	if len(verify) > 0 {
		v, err := nymeria.Verify(verify)

		if err != nil {
			fmt.Printf("Looks like an error occurred (%s).\n", err)
			return
		}

		prettyPrint(v)

		return
	}

	if len(enrich) > 0 {
		var params []nymeria.EnrichParams

		err := json.Unmarshal([]byte(enrich), &params)

		if err != nil {
			fmt.Printf("Looks like an error occurred (%s).\n", err)
			return
		}

		v, err := nymeria.Enrich(params...)

		if err != nil {
			fmt.Printf("Looks like an error occurred (%s).\n", err)
			return
		}

		prettyPrint(v)

		return
	}

	if checkAuthentication {
		if err := nymeria.CheckAuthentication(); err != nil {
			fmt.Printf("Looks like the supplied key is not valid (%s).\n", err)
			return
		}

		fmt.Println("The API key Looks good. You are ready to go!")

		return
	}

	flag.Usage()
}
