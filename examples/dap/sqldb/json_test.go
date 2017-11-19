package sqldb_test

import (
     "encoding/json"


     "github.com/influx6/moz/examples/dap"

)

var ignitorJSON = `{


    "public_id":	"",

    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0

}`

var ignitorCreateJSON = `{


    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0,

    "public_id":	"",

    "name":	""

}`

var ignitorUpdateJSON = `{


    "public_id":	"",

    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0

}`

// loadJSONFor returns a new instance of a dap.Ignitor from the provide json content.
func loadJSONFor(content string) (dap.Ignitor, error) {
	var elem dap.Ignitor

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return dap.Ignitor{}, err
	}

	return elem, nil
}

