package prompt

import (
	"fmt"

	"github.com/mielleman/nightscout-shell/internal/pkg/config"
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
	fmt.Println("prompt.Main()")
}
