package httpapi_test

import (
	"encoding/json"

	"github.com/influx6/moz/examples/dap"
)

var ignitorJSON = `{


    "public_id":	"4334343-343432-23232-332322",

    "name":	"Bob Crelly",

    "rex":	{
	
	
	    "url":	"http://github.com/wapa/api"
	
	},

    "rack":	10

}`

var ignitorCreateJSON = `{


    "public_id":	"4334343-343432-23232-332322",

    "name":	"Bob Crelly",

    "rex":	{
	
	
	    "url":	"http://github.com/wapa/api"
	
	},

    "rack":	30

}`

var ignitorUpdateJSON = `{


    "public_id":	"4334343-343432-23232-332322",

    "name":	"Bob Crelly",

    "rex":	{
	
	
	    "url":	"http://github.com/wapa/api"
	
	},

    "rack":	40

}`

// loadJSONFor returns a new instance of a dap.Ignitor from the provide json content.
func loadJSONFor(content string) (dap.Ignitor, error) {
	var elem dap.Ignitor

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return dap.Ignitor{}, err
	}

	return elem, nil
}
