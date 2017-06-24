package mongoapi_test

var ignitorJSON = `{

    "name": "",

    "public_id": "",

    "rack": 0,

    "rex": {
	
	    "url": "",
	
	},

}`

var ignitorCreateJSON = `{

    "name": "",

    "public_id": "",

    "rack": 0,

    "rex": {
	
	    "url": "",
	
	},

}`

var ignitorUpdateJSON = `{

    "name": "",

    "public_id": "",

    "rack": 0,

    "rex": {
	
	    "url": "",
	
	},

}`

// loadJSONFor returns a new instance of a dap.Ignitor from the provide json content.
func loadJSONFor(content string) (dap.Ignitor, error) {
	var elem dap.Ignitor

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return nil, err
	}

	return elem, nil
}

