
# mapstructure
    import "github.com/ardanlabs/kit/mapstructure"

Package mapstructure is a Go library for decoding generic map values to structures
and vice versa, while providing helpful error handling.

This library is most useful when decoding values from some data stream (JSON,
Gob, etc.) where you don't _quite_ know the structure of the underlying data
until you read a part of it. You can therefore read a `map[string]interface{}`
and use this library to decode it into the proper underlying native Go
structure.

Usage & Example

For usage and examples see the [Godoc](<a href="http://godoc.org/github.com/mitchellh/mapstructure">http://godoc.org/github.com/mitchellh/mapstructure</a>).

The `Decode`, `DecodePath` and `DecodeSlicePath` functions have examples associated with it there.

But Why?!

Go offers fantastic standard libraries for decoding formats such as JSON.
The standard method is to have a struct pre-created, and populate that struct
from the bytes of the encoded format. This is great, but the problem is if
you have configuration or an encoding that changes slightly depending on
specific fields. For example, consider this JSON:

	json
	{
		"type": "person",
		"name": "Mitchell"
	}

Perhaps we can't populate a specific structure without first reading
the "type" field from the JSON. We could always do two passes over the
decoding of the JSON (reading the "type" first, and the rest later).
However, it is much simpler to just decode this into a `map[string]interface{}`
structure, read the "type" key, then use something like this library
to decode it into the proper structure.

### DecodePath
Sometimes you have a large and complex JSON document where you only need to decode
a small part.


	{
		"userContext": {
			"conversationCredentials": {
		            "sessionToken": "06142010_1:75bf6a413327dd71ebe8f3f30c5a4210a9b11e93c028d6e11abfca7ff"
		    },
		    "valid": true,
		    "isPasswordExpired": false,
		    "cobrandId": 10000004,
		    "channelId": -1,
		    "locale": "en_US",
		    "tncVersion": 2,
		    "applicationId": "17CBE222A42161A3FF450E47CF4C1A00",
		    "cobrandConversationCredentials": {
		        "sessionToken": "06142010_1:b8d011fefbab8bf1753391b074ffedf9578612d676ed2b7f073b5785b"
		    },
		     "preferenceInfo": {
		         "currencyCode": "USD",
		         "timeZone": "PST",
		         "dateFormat": "MM/dd/yyyy",
		         "currencyNotationType": {
		             "currencyNotationType": "SYMBOL"
		         },
		         "numberFormat": {
		             "decimalSeparator": ".",
		             "groupingSeparator": ",",
		             "groupPattern": "###,##0.##"
		         }
		     }
		 },
		 "lastLoginTime": 1375686841,
		 "loginCount": 299,
		 "passwordRecovered": false,
		 "emailAddress": "johndoe@email.com",
		 "loginName": "sptest1",
		 "userId": 10483860,
		 "userType":
		     {
		     "userTypeId": 1,
		     "userTypeName": "normal_user"
		     }
	}

It is nice to be able to define and pull the documents and fields you need without
having to map the entire JSON structure.

	type UserType struct {
		UserTypeID   int
		UserTypeName string
	}
	
	type NumberFormat struct {
		DecimalSeparator  string `jpath:"userContext.preferenceInfo.numberFormat.decimalSeparator"`
		GroupingSeparator string `jpath:"userContext.preferenceInfo.numberFormat.groupingSeparator"`
		GroupPattern      string `jpath:"userContext.preferenceInfo.numberFormat.groupPattern"`
	}
	
	type User struct {
		Session   string   `jpath:"userContext.cobrandConversationCredentials.sessionToken"`
		CobrandID int      `jpath:"userContext.cobrandId"`
		UserType  UserType `jpath:"userType"`
		LoginName string   `jpath:"loginName"`
		NumberFormat       // This can also be a pointer to the struct (*NumberFormat)
	}
	
	doc := []byte(document)
	var docMap map[string]interface{}
	json.Unmarshal(doc, &docMap)
	
	var u User
	mapstructure.DecodePath(docMap, &u)

### DecodeSlicePath
Sometimes you have a slice of documents that you need to decode into a slice of structures


	[
		{"name":"bill"},
		{"name":"lisa"}
	]

Just Unmarshal your document into a slice of maps and decode the slice


	type NameDoc struct {
		Name string `jpath:"name"`
	}
	
	doc := []byte(document)
	var sliceMap []map[string]interface{}
	json.Unmarshal(doc, &sliceMap)
	
	var myslice []NameDoc
	err := DecodeSlicePath(sliceMap, &myslice)
	
	var myslice []*NameDoc
	err := DecodeSlicePath(sliceMap, &myslice)

### Decode Structs With Embedded Slices
Sometimes you have a document with arrays


	{
		"cobrandId": 10010352,
		"channelId": -1,
		"locale": "en_US",
		"tncVersion": 2,
		"people": [
			{
				"name": "jack",
				"age": {
				"birth":10,
				"year":2000,
				"animals": [
					{
					"barks":"yes",
					"tail":"yes"
					},
					{
					"barks":"no",
					"tail":"yes"
					}
				]
			}
			},
			{
				"name": "jill",
				"age": {
					"birth":11,
					"year":2001
				}
			}
		]
	}

You can decode within those arrays


	type Animal struct {
		Barks string `jpath:"barks"`
	}
	
	type People struct {
		Age     int      `jpath:"age.birth"` // jpath is relative to the array
		Animals []Animal `jpath:"age.animals"`
	}
	
	type Items struct {
		Categories []string `jpath:"categories"`
		Peoples    []People `jpath:"people"` // Specify the location of the array
	}
	
	docScript := []byte(document)
	var docMap map[string]interface{}
	json.Unmarshal(docScript, &docMap)
	
	var items Items
	DecodePath(docMap, &items)






## func Decode
``` go
func Decode(m interface{}, rawVal interface{}) error
```
Decode takes a map and uses reflection to convert it into the
given Go native structure. val must be a pointer to a struct.


## func DecodePath
``` go
func DecodePath(m map[string]interface{}, rawVal interface{}) error
```
DecodePath takes a map and uses reflection to convert it into the
given Go native structure. Tags are used to specify the mapping
between fields in the map and structure


## func DecodeSlicePath
``` go
func DecodeSlicePath(ms []map[string]interface{}, rawSlice interface{}) error
```
DecodeSlicePath decodes a slice of maps against a slice of structures that
contain specified tags



## type DecodeHookFunc
``` go
type DecodeHookFunc func(reflect.Kind, reflect.Kind, interface{}) (interface{}, error)
```
DecodeHookFunc declares a function to help with decoding.











## type Decoder
``` go
type Decoder struct {
    // contains filtered or unexported fields
}
```
A Decoder takes a raw interface value and turns it into structured
data, keeping track of rich error information along the way in case
anything goes wrong. Unlike the basic top-level Decode method, you can
more finely control how the Decoder behaves using the DecoderConfig
structure. The top-level Decode method is just a convenience that sets
up the most basic Decoder.









### func NewDecoder
``` go
func NewDecoder(config *DecoderConfig) (*Decoder, error)
```
NewDecoder returns a new decoder for the given configuration. Once
a decoder has been returned, the same configuration must not be used
again.


### func NewPathDecoder
``` go
func NewPathDecoder(config *DecoderConfig) (*Decoder, error)
```
NewPathDecoder returns a new decoder for the given configuration.
This is used to decode path specific structures




### func (\*Decoder) Decode
``` go
func (d *Decoder) Decode(raw interface{}) error
```
Decode decodes the given raw interface to the target pointer specified
by the configuration.



### func (\*Decoder) DecodePath
``` go
func (d *Decoder) DecodePath(m map[string]interface{}, rawVal interface{}) (bool, error)
```
DecodePath decodes the raw interface against the map based on the
specified tags



## type DecoderConfig
``` go
type DecoderConfig struct {
    // DecodeHook, if set, will be called before any decoding and any
    // type conversion (if WeaklyTypedInput is on). This lets you modify
    // the values before they're set down onto the resulting struct.
    //
    // If an error is returned, the entire decode will fail with that
    // error.
    DecodeHook DecodeHookFunc

    // If ErrorUnused is true, then it is an error for there to exist
    // keys in the original map that were unused in the decoding process
    // (extra keys).
    ErrorUnused bool

    // If WeaklyTypedInput is true, the decoder will make the following
    // "weak" conversions:
    //
    //   - bools to string (true = "1", false = "0")
    //   - numbers to string (base 10)
    //   - bools to int/uint (true = 1, false = 0)
    //   - strings to int/uint (base implied by prefix)
    //   - int to bool (true if value != 0)
    //   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
    //     FALSE, false, False. Anything else is an error)
    //   - empty array = empty map and vice versa
    //
    WeaklyTypedInput bool

    // Metadata is the struct that will contain extra metadata about
    // the decoding. If this is nil, then no metadata will be tracked.
    Metadata *Metadata

    // Result is a pointer to the struct that will contain the decoded
    // value.
    Result interface{}

    // The tag name that mapstructure reads for field names. This
    // defaults to "mapstructure"
    TagName string
}
```
DecoderConfig is the configuration that is used to create a new decoder
and allows customization of various aspects of decoding.











## type Error
``` go
type Error struct {
    Errors []string
}
```
Error implements the error interface and can represents multiple
errors that occur in the course of a single decode.











### func (\*Error) Error
``` go
func (e *Error) Error() string
```
Error is implementing the error interface.



## type Metadata
``` go
type Metadata struct {
    // Keys are the keys of the structure which were successfully decoded
    Keys []string

    // Unused is a slice of keys that were found in the raw value but
    // weren't decoded since there was no matching field in the result interface
    Unused []string
}
```
Metadata contains information about decoding a structure that
is tedious or difficult to get otherwise.

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)