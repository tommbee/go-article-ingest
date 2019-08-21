package poller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tommbee/go-article-ingest/model"
)

var configs model.Configs

// Generate a suitable poller
func Generate(URL string) (*Poller, error) {
	fn := os.Getenv("CONFIG_LOCATION")
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
			return &Poller{
				config: c,
			}, nil
		}
	}
	return &Poller{}, fmt.Errorf("Config for parser not found")
}
