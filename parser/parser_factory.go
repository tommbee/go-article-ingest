package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/tommbee/go-article-ingest/model"
)

var configs model.Configs

// Generate a suitable parser
func Generate(URL string) (*Parser, error) {
	fn := os.Getenv("CONFIG_FILE")
	jsonFile, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &configs)

	for i := 0; i < len(configs.Configs); i++ {
		c := configs.Configs[i]
		if strings.Contains(URL, c.BaseURL) {
			return &Parser{
				config: c,
			}, nil
		}
	}
	return &Parser{}, fmt.Errorf("Config for parser not found")
}
