package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	urlshort "github.com/ahmedMunna1767/url_shortener"
)

func main() {
	var urlHandler http.HandlerFunc
	inputFile := flag.String("filename", "./urlData/myData.json", "Provide a Json or Yaml file containing URL data")
	flag.Parse()

	fileType := strings.Split(*inputFile, ".")[len(strings.Split(*inputFile, "."))-1]
	mux := defaultMux()
	/* simple map for maphandler data */
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	/* map handler */
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if fileType == "yaml" {
		/* read yaml file */
		fileByteData, err := ioutil.ReadFile(*inputFile)
		if err != nil {
			panic(err)
		}
		/* yaml handler */
		urlHandler, err = urlshort.YAMLHandler([]byte(fileByteData), mapHandler)
		if err != nil {
			panic(err)
		}
	} else if fileType == "json" {
		/* read json file */
		fileByteData, err := ioutil.ReadFile(*inputFile)
		if err != nil {
			panic(err)
		}
		/* json handler */
		urlHandler, err = urlshort.JSONHandler([]byte(fileByteData), mapHandler)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Please Provide valid a yaml or json file")
		panic("")
	}

	/* start server */
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", urlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
