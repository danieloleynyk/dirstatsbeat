package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/danieloleynyk/dirstatsbeat/config"
)

// dirstatsbeat configuration.
type dirstatsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of dirstatsbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &dirstatsbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts dirstatsbeat.
func (bt *dirstatsbeat) Run(b *beat.Beat) error {
	logp.Info("dirstatsbeat is running! Hit CTRL-C to stop it.")

	var err error
	var stats os.FileInfo

	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	path := bt.config.Path
	for {
		fields := common.MapStr{}

		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		stats, err = os.Stat(path)

		if os.IsNotExist(err) {
			fields.Put("error", "Path doesn't exist")
		} else {
			fields.Put("path", path)
			fields.Put("last_updated", stats.ModTime())
			fields.Put("is_dir", stats.IsDir())
			fields.Put("size", stats.Size())
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields:    fields,
		}

		bt.client.Publish(event)
		logp.Info("Event sent")
	}
}

// Stop stops dirstatsbeat.
func (bt *dirstatsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
