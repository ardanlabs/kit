

# log
`import "github.com/ardanlabs/kit/log"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
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
	
	log.Init(os.Stderr, logLevel, log.Ldefault)

To write to the log you can make calls like this:


	log.Dev(ctx, "CreateUser", "Started : Email[%s]", nu.Email)
	log.Error(ctx, "CreateUser", err, "Completed")

The API for Dev and User follow this convention:


	log.User(ctx, "funcName", "formatted message %s", values)

ctx

A ctx is a unique id that can be used to trace an entire session or
request. This value should be generated as early as possible and passed
through out the different calls.

funcName

Provide the name of the function the log statement is being declared in. This
can take on a type name in the case of the method (type.method).

formatted message, values

Any string can be provided but it does support a formatted message. Values
would be substituted if provided. This messaging is up to you.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func Dev(traceID string, funcName string, format string, a ...interface{})](#Dev)
* [func DevOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#DevOffset)
* [func Error(traceID string, funcName string, err error, format string, a ...interface{})](#Error)
* [func ErrorOffset(traceID string, offset int, funcName string, err error, format string, a ...interface{})](#ErrorOffset)
* [func Fatal(traceID string, funcName string, format string, a ...interface{})](#Fatal)
* [func FatalOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#FatalOffset)
* [func Init(w io.Writer, level func() int, flags int)](#Init)
* [func User(traceID string, funcName string, format string, a ...interface{})](#User)
* [func UserOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#UserOffset)
* [type Logger](#Logger)
  * [func New(w io.Writer, levelHandler func() int, flags int) *Logger](#New)
  * [func (l *Logger) Dev(traceID string, funcName string, format string, a ...interface{})](#Logger.Dev)
  * [func (l *Logger) DevOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#Logger.DevOffset)
  * [func (l *Logger) Error(traceID string, funcName string, err error, format string, a ...interface{})](#Logger.Error)
  * [func (l *Logger) ErrorOffset(traceID string, offset int, funcName string, err error, format string, a ...interface{})](#Logger.ErrorOffset)
  * [func (l *Logger) Fatal(traceID string, funcName string, format string, a ...interface{})](#Logger.Fatal)
  * [func (l *Logger) FatalOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#Logger.FatalOffset)
  * [func (l *Logger) User(traceID string, funcName string, format string, a ...interface{})](#Logger.User)
  * [func (l *Logger) UserOffset(traceID string, offset int, funcName string, format string, a ...interface{})](#Logger.UserOffset)

#### <a name="pkg-examples">Examples</a>
* [Dev](#example_Dev)

#### <a name="pkg-files">Package files</a>
[doc.go](/src/github.com/ardanlabs/kit/log/doc.go) [log.go](/src/github.com/ardanlabs/kit/log/log.go) [log_default.go](/src/github.com/ardanlabs/kit/log/log_default.go) 


## <a name="pkg-constants">Constants</a>
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



## <a name="Dev">func</a> [Dev](/src/target/log_default.go?s=405:479#L12)
``` go
func Dev(traceID string, funcName string, format string, a ...interface{})
```
Dev logs trace information for developers.



## <a name="DevOffset">func</a> [DevOffset](/src/target/log_default.go?s=1207:1299#L33)
``` go
func DevOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
DevOffset logs trace information for developers with a offset option to
expand the caller level.



## <a name="Error">func</a> [Error](/src/target/log_default.go?s=756:843#L22)
``` go
func Error(traceID string, funcName string, err error, format string, a ...interface{})
```
Error logs trace information that are errors.



## <a name="ErrorOffset">func</a> [ErrorOffset](/src/target/log_default.go?s=1722:1827#L45)
``` go
func ErrorOffset(traceID string, offset int, funcName string, err error, format string, a ...interface{})
```
ErrorOffset logs trace information that are errors with a offset option to
expand the caller level.



## <a name="Fatal">func</a> [Fatal](/src/target/log_default.go?s=971:1047#L27)
``` go
func Fatal(traceID string, funcName string, format string, a ...interface{})
```
Fatal logs trace information for users and terminates the app.



## <a name="FatalOffset">func</a> [FatalOffset](/src/target/log_default.go?s=2009:2103#L51)
``` go
func FatalOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
FatalOffset logs trace information for users and terminates the app with a
offset expand the caller level.



## <a name="Init">func</a> [Init](/src/target/log_default.go?s=194:245#L1)
``` go
func Init(w io.Writer, level func() int, flags int)
```
Init initializes the default logger to allow usage of the global log
functions.



## <a name="User">func</a> [User](/src/target/log_default.go?s=576:651#L17)
``` go
func User(traceID string, funcName string, format string, a ...interface{})
```
User logs trace information for users.



## <a name="UserOffset">func</a> [UserOffset](/src/target/log_default.go?s=1460:1553#L39)
``` go
func UserOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
UserOffset logs trace information for users with a offset option to expand the
caller level.




## <a name="Logger">type</a> [Logger](/src/target/log.go?s=946:1019#L28)
``` go
type Logger struct {
    *log.Logger
    // contains filtered or unexported fields
}
```
Logger contains a standard logger for all logging.







### <a name="New">func</a> [New](/src/target/log.go?s=1060:1125#L35)
``` go
func New(w io.Writer, levelHandler func() int, flags int) *Logger
```
New returns a instance of a logger.





### <a name="Logger.Dev">func</a> (\*Logger) [Dev](/src/target/log.go?s=1339:1425#L46)
``` go
func (l *Logger) Dev(traceID string, funcName string, format string, a ...interface{})
```
Dev logs trace information for developers.




### <a name="Logger.DevOffset">func</a> (\*Logger) [DevOffset](/src/target/log.go?s=2844:2948#L109)
``` go
func (l *Logger) DevOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
DevOffset logs trace information for developers with a offset option to
expand the caller level.




### <a name="Logger.Error">func</a> (\*Logger) [Error](/src/target/log.go?s=2032:2131#L76)
``` go
func (l *Logger) Error(traceID string, funcName string, err error, format string, a ...interface{})
```
Error logs trace information that are errors.




### <a name="Logger.ErrorOffset">func</a> (\*Logger) [ErrorOffset](/src/target/log.go?s=3701:3818#L141)
``` go
func (l *Logger) ErrorOffset(traceID string, offset int, funcName string, err error, format string, a ...interface{})
```
ErrorOffset logs trace information that are errors with a offset option to
expand the caller level.




### <a name="Logger.Fatal">func</a> (\*Logger) [Fatal](/src/target/log.go?s=2424:2512#L91)
``` go
func (l *Logger) Fatal(traceID string, funcName string, format string, a ...interface{})
```
Fatal logs trace information for users and terminates the app.




### <a name="Logger.FatalOffset">func</a> (\*Logger) [FatalOffset](/src/target/log.go?s=4165:4271#L157)
``` go
func (l *Logger) FatalOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
FatalOffset logs trace information for users and terminates the app with a
offset expand the caller level.




### <a name="Logger.User">func</a> (\*Logger) [User](/src/target/log.go?s=1681:1768#L61)
``` go
func (l *Logger) User(traceID string, funcName string, format string, a ...interface{})
```
User logs trace information for users.




### <a name="Logger.UserOffset">func</a> (\*Logger) [UserOffset](/src/target/log.go?s=3268:3373#L125)
``` go
func (l *Logger) UserOffset(traceID string, offset int, funcName string, format string, a ...interface{})
```
UserOffset logs trace information for users with a offset option to expand the
caller level.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
