// Sample program to show how to use the mapstructure package to handle large
// JSON documents and flatten them out.
package main

import (
	"encoding/json"
	"fmt"

	ms "github.com/ardanlabs/kit/mapstructure"
)

func main() {
	decodePath()
	decodeSlice()
	decodePathArray()
}

//==============================================================================

// decodePath shows how to use the DecodePath API call.
func decodePath() {
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
		Session      string   `jpath:"userContext.cobrandConversationCredentials.sessionToken"`
		CobrandID    int      `jpath:"userContext.cobrandId"`
		UserType     UserType `jpath:"userType"`
		LoginName    string   `jpath:"loginName"`
		NumberFormat          // This can also be a pointer to the struct (*NumberFormat).
	}

	doc := []byte(decodePathDoc)
	var docMap map[string]interface{}
	json.Unmarshal(doc, &docMap)

	var u User
	ms.DecodePath(docMap, &u)

	fmt.Printf("DecodePath : %+v\n", u)
}

var decodePathDoc = `{
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
}`

//==============================================================================

// decodeSlice shows how to use the DecodeSlicePath API call.
func decodeSlice() {
	type NameDoc struct {
		Name string `jpath:"name"`
	}

	doc := []byte(decodeSliceDoc)
	var sliceMap []map[string]interface{}
	json.Unmarshal(doc, &sliceMap)

	var byValue []NameDoc
	ms.DecodeSlicePath(sliceMap, &byValue)

	fmt.Printf("DecodeSlicePath : ByValue : %+v\n", byValue)

	var byAddr []*NameDoc
	ms.DecodeSlicePath(sliceMap, &byAddr)

	fmt.Printf("DecodeSlicePath : ByAddr : %+v\n", byAddr)
}

var decodeSliceDoc = `[
	{"name":"bill"},
	{"name":"lisa"}
]`

//==============================================================================

// decodePath shows how to use the DecodePath API call with an array of
// sub documents.
func decodePathArray() {
	type Animal struct {
		Barks string `jpath:"barks"`
	}

	type People struct {
		Name    string   `jpath:"name"` // jpath is relative to the array.
		Age     int      `jpath:"age.birth"`
		Animals []Animal `jpath:"age.animals"`
	}

	type Items struct {
		Peoples []People `jpath:"people"` // Specify the location of the array.
	}

	doc := []byte(decodePathArrayDoc)
	var docMap map[string]interface{}
	json.Unmarshal(doc, &docMap)

	var items Items
	ms.DecodePath(docMap, &items)

	fmt.Printf("DecodePathArray : %+v\n", items)
}

var decodePathArrayDoc = `{
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
}`
