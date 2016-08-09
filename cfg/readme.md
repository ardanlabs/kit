
# cfg
    import "github.com/ardanlabs/kit/cfg"

Package cfg provides configuration options that are loaded from the environment.
Configuration is then stored in memory and can be retrieved by its proper
type.

To initialize the configuration system from your environment, call Init:


	cfg.Init(cfg.EnvProvider{Namespace: "configKey"})

To retrieve values from configuration:


	proc, err := cfg.String("proc_id")
	port, err := cfg.Int("port")
	ms, err := cfg.Time("stamp")

Use the Must set of function to retrieve a single value but these calls
will panic if the key does not exist:


	proc := cfg.MustString("proc_id")
	port := cfg.MustInt("port")
	ms := cfg.MustTime("stamp")






## func Bool
``` go
func Bool(key string) (bool, error)
```
Bool calls the default Config and returns the bool value of a given key as a
bool. It will return an error if the key was not found or the value can't be
converted to a bool.


## func Duration
``` go
func Duration(key string) (time.Duration, error)
```

## func Init
``` go
func Init(p Provider) error
```
Init populates the package's default Config and should be called only once.
A Provider must be supplied which will return a map of key/value pairs to be
loaded.


## func Int
``` go
func Int(key string) (int, error)
```
Int calls the Default config and returns the value of the given key as an
int. It will return an error if the key was not found or the value
can't be converted to an int.


## func Log
``` go
func Log() string
```
Log returns a string to help with logging the package's default Config. It
excludes any values whose key contains the string "PASS".


## func MustBool
``` go
func MustBool(key string) bool
```
MustBool calls the default Config and returns the bool value of a given key
as a bool. It will panic if the key was not found or the value can't be
converted to a bool.


## func MustDuration
``` go
func MustDuration(key string) time.Duration
```

## func MustInt
``` go
func MustInt(key string) int
```
MustInt calls the default Config and returns the value of the given key as
an int. It will panic if the key was not found or the value can't be
converted to an int.


## func MustString
``` go
func MustString(key string) string
```
MustString calls the default Config and returns the value of the given key
as a string, else it will panic if the key was not found.


## func MustTime
``` go
func MustTime(key string) time.Time
```
MustTime calls the default Config ang returns the value of the given key as
a Time. It will panic if the key was not found or the value can't be
converted to a Time.


## func MustURL
``` go
func MustURL(key string) *url.URL
```
MustURL calls the default Config and returns the value of the given key as a
URL. It will panic if the key was not found or the value can't be converted
to a URL.


## func SetBool
``` go
func SetBool(key string, value bool)
```
SetBool adds or modifies the default Config for the specified key and value.


## func SetDuration
``` go
func SetDuration(key string, value time.Duration)
```

## func SetInt
``` go
func SetInt(key string, value int)
```
SetInt adds or modifies the default Config for the specified key and value.


## func SetString
``` go
func SetString(key string, value string)
```
SetString adds or modifies the default Config for the specified key and
value.


## func SetTime
``` go
func SetTime(key string, value time.Time)
```
SetTime adds or modifies the default Config for the specified key and value.


## func SetURL
``` go
func SetURL(key string, value *url.URL)
```
SetURL adds or modifies the default Config for the specified key and value.


## func String
``` go
func String(key string) (string, error)
```
String calls the default Config and returns the value of the given key as a
string. It will return an error if key was not found.


## func Time
``` go
func Time(key string) (time.Time, error)
```
Time calls the default Config and returns the value of the given key as a
Time. It will return an error if the key was not found or the value can't be
converted to a Time.


## func URL
``` go
func URL(key string) (*url.URL, error)
```
URL calls the default Config and returns the value of the given key as a
URL. It will return an error if the key was not found or the value can't be
converted to a URL.



## type Config
``` go
type Config struct {
    // contains filtered or unexported fields
}
```
Config is a goroutine safe configuration store, with a map of values
set from a config Provider.









### func New
``` go
func New(p Provider) (*Config, error)
```
New populates a new Config from a Provider. It will return an error if there
was any problem reading from the Provider.




### func (\*Config) Bool
``` go
func (c *Config) Bool(key string) (bool, error)
```
Bool returns the bool value of a given key as a bool. It will return an
error if the key was not found or the value can't be converted to a bool.



### func (\*Config) Duration
``` go
func (c *Config) Duration(key string) (time.Duration, error)
```
Duration returns the value of the given key as a Duration. It will return an
error if the key was not found or the value can't be converted to a Duration.



### func (\*Config) Int
``` go
func (c *Config) Int(key string) (int, error)
```
Int returns the value of the given key as an int. It will return an error if
the key was not found or the value can't be converted to an int.



### func (\*Config) Log
``` go
func (c *Config) Log() string
```
Log returns a string to help with logging your configuration. It excludes
any values whose key contains the string "PASS".



### func (\*Config) MustBool
``` go
func (c *Config) MustBool(key string) bool
```
MustBool returns the bool value of a given key as a bool. It will panic if
the key was not found or the value can't be converted to a bool.



### func (\*Config) MustDuration
``` go
func (c *Config) MustDuration(key string) time.Duration
```
MustDuration returns the value of the given key as a Duration. It will panic
if the key was not found or the value can't be converted into a Duration.



### func (\*Config) MustInt
``` go
func (c *Config) MustInt(key string) int
```
MustInt returns the value of the given key as an int. It will panic if the
key was not found or the value can't be converted to an int.



### func (\*Config) MustString
``` go
func (c *Config) MustString(key string) string
```
MustString returns the value of the given key as a string. It will panic if
the key was not found.



### func (\*Config) MustTime
``` go
func (c *Config) MustTime(key string) time.Time
```
MustTime returns the value of the given key as a Time. It will panic if the
key was not found or the value can't be converted to a Time.



### func (\*Config) MustURL
``` go
func (c *Config) MustURL(key string) *url.URL
```
MustURL returns the value of the given key as a URL. It will panic if the
key was not found or the value can't be converted to a URL.



### func (\*Config) SetBool
``` go
func (c *Config) SetBool(key string, value bool)
```
SetBool adds or modifies the configuration for the specified key and value.



### func (\*Config) SetDuration
``` go
func (c *Config) SetDuration(key string, value time.Duration)
```
SetDuration adds or modifies the configuration for a given duration at a
specific key.



### func (\*Config) SetInt
``` go
func (c *Config) SetInt(key string, value int)
```
SetInt adds or modifies the configuration for the specified key and value.



### func (\*Config) SetString
``` go
func (c *Config) SetString(key string, value string)
```
SetString adds or modifies the configuration for the specified key and
value.



### func (\*Config) SetTime
``` go
func (c *Config) SetTime(key string, value time.Time)
```
SetTime adds or modifies the configuration for the specified key and value.



### func (\*Config) SetURL
``` go
func (c *Config) SetURL(key string, value *url.URL)
```
SetURL adds or modifies the configuration for the specified key and value.



### func (\*Config) String
``` go
func (c *Config) String(key string) (string, error)
```
String returns the value of the given key as a string. It will return an
error if key was not found.



### func (\*Config) Time
``` go
func (c *Config) Time(key string) (time.Time, error)
```
Time returns the value of the given key as a Time. It will return an error
if the key was not found or the value can't be converted to a Time.



### func (\*Config) URL
``` go
func (c *Config) URL(key string) (*url.URL, error)
```
URL returns the value of the given key as a URL. It will return an error if
the key was not found or the value can't be converted to a URL.



## type EnvProvider
``` go
type EnvProvider struct {
    Namespace string
}
```
EnvProvider provides configuration from the environment. All keys will be
made uppercase.











### func (EnvProvider) Provide
``` go
func (ep EnvProvider) Provide() (map[string]string, error)
```
Provide implements the Provider interface.



## type FileProvider
``` go
type FileProvider struct {
    Filename string
}
```
FileProvider describes a file based loader which loads the configuration
from a file listed.











### func (FileProvider) Provide
``` go
func (fp FileProvider) Provide() (map[string]string, error)
```
Provide implements the Provider interface.



## type MapProvider
``` go
type MapProvider struct {
    Map map[string]string
}
```
MapProvider provides a simple implementation of the Provider whereby it just
returns a stored map.











### func (MapProvider) Provide
``` go
func (mp MapProvider) Provide() (map[string]string, error)
```
Provide implements the Provider interface.



## type Provider
``` go
type Provider interface {
    Provide() (map[string]string, error)
}
```
Provider is implemented by the user to provide the configuration as a map.
There are currently two Providers implemented, EnvProvider and MapProvider.

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)