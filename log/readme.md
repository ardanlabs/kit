
# log
    import "github.com/ardanlabs/kit/log"

Package log provides a simple layer above the standard library logging package.
The base API provides two logging levels, DEV and USER. The DEV level logs things
developers need and can be verbose. The USER level logs things for users need
and should not be verbose. There is an Error call which falls under USER.

To initialize the logging system from your application, call Init:


	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}
	
	log.Init(os.Stderr, logLevel)

To write to the log you can make calls like this:


	log.Dev(context, "CreateUser", "Started : Email[%s]", nu.Email)
	log.Error(context, "CreateUser", err, "Completed")

The API for Dev and User follow this convention:


	log.User(context, "funcName", "formatted message %s", values)

context

A context is a unique id that can be used to trace an entire session or
request. This value should be generated as early as possible and passed
through out the different calls.

funcName

Provide the name of the function the log statement is being declared in. This
can take on a type name in the case of the method (type.method).

formatted message, values

Any string can be provided but it does support a formatted message. Values
would be substituted if provided. This messaging is up to you.




## Constants
``` go
const (
    NONE int = iota
    DEV
    USER
)
```
Level constants that define the supported usable LogLevel.

``` go
const (
    // Ldate enables the date in the local time zone: 2009/01/23
    Ldate = 1 << iota
    // Ltime enables the time in the local time zone: 01:23:23
    Ltime
    // Lmicroseconds enables microsecond resolution: 01:23:23.123123.  assumes Ltime.
    Lmicroseconds
    // Llongfile enables full file name and line number: /a/b/c/d.go:23
    Llongfile
    // Lshortfile enables final file name element and line number: d.go:23. overrides Llongfile
    Lshortfile
    // LUTC enables if Ldate or Ltime is set, use UTC rather than the local time zone
    LUTC
    // LstdFlags enables initial values for the standard logger
    LstdFlags = Ldate | Ltime
    // Ldefault enables intial values for the default kit logger
    Ldefault = log.Ldate | log.Ltime | log.Lshortfile
)
```


## func Dev
``` go
func Dev(context interface{}, funcName string, format string, a ...interface{})
```
Dev logs trace information for developers.


## func DevOffset
``` go
func DevOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
DevOffset logs trace information for developers with a offset option to
expand the caller level.


## func Error
``` go
func Error(context interface{}, funcName string, err error, format string, a ...interface{})
```
Error logs trace information that are errors.


## func ErrorOffset
``` go
func ErrorOffset(context interface{}, offset int, funcName string, err error, format string, a ...interface{})
```
ErrorOffset logs trace information that are errors with a offset option to
expand the caller level.


## func Fatal
``` go
func Fatal(context interface{}, funcName string, format string, a ...interface{})
```
Fatal logs trace information for users and terminates the app.


## func FatalOffset
``` go
func FatalOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
FatalOffset logs trace information for users and terminates the app with a
offset expand the caller level.


## func Init
``` go
func Init(w io.Writer, level func() int, flags int)
```
Init initializes the default logger to allow usage of the global log
functions.


## func User
``` go
func User(context interface{}, funcName string, format string, a ...interface{})
```
User logs trace information for users.


## func UserOffset
``` go
func UserOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
UserOffset logs trace information for users with a offset option to expand the
caller level.



## type Logger
``` go
type Logger struct {
    *log.Logger
    // contains filtered or unexported fields
}
```
Logger contains a standard logger for all logging.









### func New
``` go
func New(w io.Writer, levelHandler func() int, flags int) *Logger
```
New returns a instance of a logger.




### func (\*Logger) Dev
``` go
func (l *Logger) Dev(context interface{}, funcName string, format string, a ...interface{})
```
Dev logs trace information for developers.



### func (\*Logger) DevOffset
``` go
func (l *Logger) DevOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
DevOffset logs trace information for developers with a offset option to
expand the caller level.



### func (\*Logger) Error
``` go
func (l *Logger) Error(context interface{}, funcName string, err error, format string, a ...interface{})
```
Error logs trace information that are errors.



### func (\*Logger) ErrorOffset
``` go
func (l *Logger) ErrorOffset(context interface{}, offset int, funcName string, err error, format string, a ...interface{})
```
ErrorOffset logs trace information that are errors with a offset option to
expand the caller level.



### func (\*Logger) Fatal
``` go
func (l *Logger) Fatal(context interface{}, funcName string, format string, a ...interface{})
```
Fatal logs trace information for users and terminates the app.



### func (\*Logger) FatalOffset
``` go
func (l *Logger) FatalOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
FatalOffset logs trace information for users and terminates the app with a
offset expand the caller level.



### func (\*Logger) User
``` go
func (l *Logger) User(context interface{}, funcName string, format string, a ...interface{})
```
User logs trace information for users.



### func (\*Logger) UserOffset
``` go
func (l *Logger) UserOffset(context interface{}, offset int, funcName string, format string, a ...interface{})
```
UserOffset logs trace information for users with a offset option to expand the
caller level.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)