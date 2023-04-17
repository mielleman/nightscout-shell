package service

import (
	"fmt"
	"os"
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
		log.Error("ERROR: Failed %s", err)
		s.Stop()
	}

	contents := s.convertEntry(entry)
	fmt.Println(contents)

	// Save it to the cache file
	err = s.writeCache(contents)
	if err != nil {
		// Every error is fatal
		log.Error("ERROR: Failed %s", err)
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
	// (work with the raw value, as it van be compared to the thresholds)
	var color string
	if entry.Value >= s.nightscout.Status.Settings.Thresholds.High {
		color = "HIGH"
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.TargetTop {
		color = "ABOVE"
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.TargetBottom {
		color = "OK"
	} else if entry.Value >= s.nightscout.Status.Settings.Thresholds.Low {
		color = "BELOW"
	} else {
		color = "LOW"
	}

	// Output the string
	return fmt.Sprintf("%s%s (%s)", value, direction, color)
}

func (s *Service) writeCache(contents string) error {
	file, err := os.Create(s.config.CacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
