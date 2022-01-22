package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYamlData, err := myParseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := myBuildMap(parsedYamlData)
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback mapHandler will be called instead.
//
// JSON is expected to be in the format:
//
//  {
//      "path": "/urlshort",
//      "url": "https://github.com/gophercises/urlshort"
//  }
//

func JSONHandler(JsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJsonData, err := myParsejson(JsonBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := myBuildMap(parsedJsonData)
	return MapHandler(pathsToUrls, fallback), nil
}

func myBuildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}
	return pathToUrls
}

func myParseYaml(yamlData []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yamlData, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func myParsejson(jsonData []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := json.Unmarshal(jsonData, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
