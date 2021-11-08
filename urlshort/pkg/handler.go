package urlshortener

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	fmt.Println(pathsToUrls)
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type PathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

type convert func([]byte, interface{}) error

func toMap(pu []PathURL) map[string]string {
	pathMap := make(map[string]string)
	for _, v := range pu {
		pathMap[v.Path] = v.URL
	}
	return pathMap
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
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return handle(yml, fallback, yaml.Unmarshal)
}

func JSONHandler(jdata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return handle(jdata, fallback, json.Unmarshal)
}

func handle(data []byte, fallback http.Handler, fn convert) (http.HandlerFunc, error) {
	out := []PathURL{}
	err := fn(data, &out)
	if err != nil {
		return nil, err
	}
	pathMap := toMap(out)
	return MapHandler(pathMap, fallback), nil
}

func BoltDBHandler(db *bolt.DB, fallback http.Handler, bName string) (http.HandlerFunc, error) {
	paths := []PathURL{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			p := PathURL{}
			json.Unmarshal(v, &p)
			paths = append(paths, p)
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not read bucket %s: %v", bName, err)
	}
	pathMap := toMap(paths)
	return MapHandler(pathMap, fallback), nil
}
