package config

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"
)

const (
	KeyDeviceUUID = "net.nitroshare.device_uuid"
	KeyDeviceName = "net.nitroshare.device_name"
)

type configValue struct {
	value    string
	watchers []chan string
}

func (c configValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c *configValue) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.value)
}

// Config stores the local configuration for the application and provides
// thread-safe access to it.
type Config struct {
	mutex      sync.Mutex
	filename   string
	values     map[string]*configValue
	chanSync   chan any
	chanClosed chan any
}

func (c *Config) sync() error {
	f, err := os.Create(c.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return json.NewEncoder(f).Encode(c.values)
}

func (c *Config) run() {
	defer func() {
		for _, v := range c.values {
			for _, w := range v.watchers {
				close(w)
			}
		}
		close(c.chanClosed)
	}()
	var chanT <-chan time.Time
	for {
		select {
		case <-chanT:
			if err := c.sync(); err != nil {
				// TODO: log
			}
		case _, ok := <-c.chanSync:
			if !ok {
				return
			}
			chanT = time.After(500 * time.Millisecond)
		}
	}
}

// New creates a new Config instance and loads existing values.
func New(configPath string) (*Config, error) {
	var (
		filename = path.Join(configPath, "config.json")
		values   = map[string]*configValue{}
	)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&values); err != nil {
		return nil, err
	}
	return &Config{
		filename:   filename,
		values:     values,
		chanSync:   make(chan any),
		chanClosed: make(chan any),
	}, nil
}

func (c *Config) set(key string, v *configValue) {
	c.values[key] = v
	for _, w := range v.watchers {
		w <- v.value
	}
	c.chanSync <- nil
}

// Get retrieves the current value for the provided key. If no value is
// currently set, the provided default is set and then returned.
func (c *Config) Get(key, def string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, ok := c.values[key]
	if !ok {
		v = &configValue{value: def}
		c.set(key, v)
	}
	return v.value
}

// Set changes the current value for the provided key. Any watchers set on the
// key will be sent the new value. Changes are written to disk shortly after.
func (c *Config) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, ok := c.values[key]
	if !ok {
		v = &configValue{}
	}
	if value == v.value {
		return
	}
	v.value = value
	c.set(key, v)
}

// Watch returns a channel that can be used to read changes to the provided value.
func (c *Config) Watch(key string) <-chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	chanWatch := make(chan string)
	v, ok := c.values[key]
	if !ok {
		v = &configValue{}
		c.values[key] = v
	}
	v.watchers = append(v.watchers, chanWatch)
	return chanWatch
}

// Close shuts down the config and all watchers. No other methods should be
// called after this method.
func (c *Config) Close() {
	close(c.chanSync)
	<-c.chanClosed
}
