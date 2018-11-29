package normaliser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tommbee/go-article-ingest/model"
)

var configs model.Configs

// Generate a suitable normaliser
func Generate(URL string) (*Normaliser, error) {
	fn := os.Getenv("CONFIG_FILE")
	jsonFile, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &configs)

	for i := 0; i < len(configs.Configs); i++ {
		c := configs.Configs[i]
		if strings.Contains(URL, c.BaseURL) {
			return &Normaliser{
				DateFormat: c.DateFormat,
			}, nil
		}
	}
	return &Normaliser{}, fmt.Errorf("Config for normaliser not found")
}
