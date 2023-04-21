package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mielleman/nightscout-shell/internal/pkg/config"
	"github.com/mielleman/nightscout-shell/internal/pkg/nightscout"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	config     *config.Configuration
	stop       chan int
	ticker     *time.Ticker
	nightscout *nightscout.Nightscout
}

func New(filename string) *Service {
	s := &Service{
		config: config.ParseConfigFile(filename),
		stop:   make(chan int),
	}

	// Create the Nightscout instance (from the config)
	s.nightscout = nightscout.New(s.config.NightscoutUrl, s.config.NightscoutToken)

	// Get the latest status from Nightscout
	err := s.nightscout.GetStatus()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return s
}

func (s *Service) Start() {
	// Create the ticker
	s.ticker = time.NewTicker(time.Duration(s.config.ServiceInterval) * time.Minute)
	s.run()

	// Now 'hang'
	<-s.stop
}

func (s *Service) Stop() {
	// Stop the routine
	close(s.stop)
}

func (s *Service) run() {
	go func() {
		// We always run upon start
		s.update()

		// Now keep ticking until the stop channel signals any value
		for {
			select {
			case <-s.ticker.C:
				s.update()
			case <-s.stop:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Service) update() {
	// Retrieve the latest value
	entry, err := s.nightscout.GetLastEntry()
	if err != nil {
		// Every error is fatal
		log.Errorf("ERROR: Failed %s", err)
		s.Stop()
	}

	contents := s.convertEntry(entry)
	log.WithField("prompt", contents).Infof("Received %d mg/dl, converting to prompt", entry.Value)

	// Save it to the cache file
	err = s.writeCache(contents)
	if err != nil {
		// Every error is fatal
		log.Errorf("ERROR: Failed %s", err)
		s.Stop()
	}
}

func (s *Service) convertEntry(entry *nightscout.Entry) string {
	// Determine the value
	var value string
	if s.nightscout.Status.Settings.Units == "mmol" {
		// convert the mg/dl to mmol/l
		value = fmt.Sprintf("%.1f", float64(entry.Value)/18)
	} else {
		// just use the mg/dl
		value = fmt.Sprintf("%d", entry.Value)
	}

	// Determine the direction symbol
	var direction string
	switch entry.Direction {
	case "DoubleUp":
		direction = "⇈"
	case "SingleUp":
		direction = "↑"
	case "FortyFiveUp":
		direction = "↗"
	case "Flat":
		direction = "→"
	case "FortyFiveDown":
		direction = "↘"
	case "SingleDown":
		direction = "↓"
	case "DoubleDown":
		direction = "⇊"
	default:
		direction = "?"
	}

	// Check the Thresholds
	// (work with the raw mgdl value, as it then can be compared to the thresholds)
	var color string
	if entry.Value >= s.nightscout.Status.Settings.Thresholds.High {
		color = "1;31" // High
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.TargetTop {
		color = "0;33" // Above
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.TargetBottom {
		color = "0;32" // Ok
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.Low {
		color = "0;33" // Below
	} else {
		color = "1;31" // Low
	}

	log.WithFields(log.Fields{
		"color":     color,
		"value":     value,
		"direction": direction,
	}).Debug("Outcome")

	// Output the string
	return fmt.Sprintf("\x1b[%sm💉 %s %s\x1b[m", color, value, direction)
}

func (s *Service) writeCache(contents string) error {
	// Create the caching directory if it does not exist
	err := os.MkdirAll(filepath.Dir(s.config.CacheFile), 0700)
	if err != nil {
		return err
	}

	// Open/Create the cache file
	file, err := os.Create(s.config.CacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content
	_, err = file.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
