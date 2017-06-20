package mongoapi_test

import (
     "testing"


     "encoding/json"


     "github.com/influx6/faux/metrics/sentries/stdout"


     "github.com/influx6/moz/examples/dap"


     "github.com/influx6/moz/examples/dap/mongoapi"

)

var IgnitorJSON = `{
    {
	"name": "\"\"",
	"public_id": "\"\"",
	"rack": "0",
	"rex": "{\n\t\"url\": \"\\\"\\\"\"\n}"
}
}`

var IgnitorCreateJSON = `{
    {
	"name": "\"\"",
	"public_id": "\"\"",
	"rack": "0",
	"rex": "{\n\t\"url\": \"\\\"\\\"\"\n}"
}
}`

var IgnitorUpdateJSON = `{
    {
	"name": "\"\"",
	"public_id": "\"\"",
	"rack": "0",
	"rex": "{\n\t\"url\": \"\\\"\\\"\"\n}"
}
}`