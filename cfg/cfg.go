package cfg

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// c represents the configuration store, with a map to store the loaded keys
// from the environment.
var c struct {
	m  map[string]string
	mu sync.RWMutex
}

//==============================================================================

// Provider is implemented by the user to provide the configuration as a map
type Provider interface {
	Provide() (map[string]string, error)
}

//==============================================================================

// Init is to be called only once, to load up the given namespace if found,
// in the environment variables. All keys will be made lowercase.
func Init(p Provider) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get the provided configuration.
	m, err := p.Provide()
	if err != nil {
		return err
	}

	// Set it to the global instance.
	c.m = m

	return nil
}

// Log returns a string to help with logging configuration.
func Log() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var buf bytes.Buffer
	for k, v := range c.m {
		if !strings.Contains(k, "PASS") {
			buf.WriteString(k + "=" + v + "\n")
		}
	}

	return buf.String()
}

// String returns the value of the given key as a string, else it will return
// an error if key was not found.
func String(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		return "", fmt.Errorf("Unknown key %s !", key)
	}

	return value, nil
}

// MustString returns the value of the given key as a string, else it will panic
// if the key was not found.
func MustString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown key %s !", key))
	}

	return value
}

// Int returns the value of the given key as an int, else it will return
// an error, if the key was not found or the value can't be convered to an int.
func Int(key string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		return 0, fmt.Errorf("Unknown Key %s !", key)
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return iv, nil
}

// MustInt returns the value of the given key as an int, else it will panic
// if the key was not found or the value can't be convered to an int.
func MustInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Key %q value is not an int", key))
	}

	return iv
}

// Time returns the value of the given key as a Time, else it will return an
// error, if the key was not found or the value can't be convered to a Time.
func Time(key string) (time.Time, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		return time.Time{}, fmt.Errorf("Unknown Key %s !", key)
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		return tv, err
	}

	return tv, nil
}

// MustTime returns the value of the given key as a Time, else it will panic
// if the key was not found or the value can't be convered to a Time.
func MustTime(key string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		panic(fmt.Sprintf("Key %q value is not a Time", key))
	}

	return tv
}

// Bool returns the bool balue of a given key as a bool, else it will return an
// error, if the key was not found or the value can't be convered to a bool.
func Bool(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		return false, fmt.Errorf("Unknown Key %s !", key)
	}

	if value == "on" || value == "yes" {
		value = "true"
	} else if value == "off" || value == "no" {
		value = "false"
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return val, nil
}

// MustBool returns the bool balue of a given key as a bool, else it will panic
// if the key was not found or the value can't be convered to a bool.
func MustBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	if value == "on" || value == "yes" {
		value = "true"
	} else if value == "off" || value == "no" {
		value = "false"
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return val
}

// URL returns the value of the given key as a URL, else it will return an
// error, if the key was not found or the value can't be convered to a URL.
func URL(key string) (*url.URL, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		return nil, fmt.Errorf("Unknown Key %s !", key)
	}

	u, err := url.Parse(value)
	if err != nil {
		return u, err
	}

	return u, nil
}

// MustURL returns the value of the given key as a URL, else it will panic
// if the key was not found or the value can't be convered to a URL.
func MustURL(key string) *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	u, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("Key %q value is not a URL", key))
	}

	return u
}
