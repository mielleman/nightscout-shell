package prompt

import (
	"fmt"
	"os"

	"github.com/mielleman/nightscout-shell/internal/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Prompt struct {
	config *config.Configuration
}

func New(filename string) *Prompt {
	return &Prompt{
		config: config.ParseConfigFile(filename),
	}
}

func (p *Prompt) Main() {
	data, err := os.ReadFile(p.config.CacheFile)
	if err != nil {
		log.Errorf("Could not read the cache file '%s'", p.config.CacheFile)
		log.Error(err)
		os.Exit(2)
	}
	fmt.Printf("%s", data)
}
