Ignitor HTTP API 
===============================

Ignitor HTTP API attempts to provide a simple basic documentation which details 
the basic structure of the Ignitor type, and the response which will be provided 
when working with this API.

The API provides a basic CRUD http API as describe below:

## Create { POST /Ignitor/ }

Create receives the provided record of the Ignitor type which is delieved the 
JSON content to the HTTP API. This will in turn return a respective status code.

- Expected Content Type: 

```http
    Content-Type: application/json
```

- Expected Request Parameters

```json
    {
    }
```

- Expected Request Body

```json
{


    "public_id":	"",

    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0

}
```

- Expected Status Code

```
Failure: 500
Success: 201
```

- Expected Response Body

```json
{


    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0,

    "public_id":	""

}
```

## GET /Ignitor/:public_id

Get retrives a giving record of the Ignitor type from the HTTP API returning received result as a JSON
response. It uses the provided `:public_id` parameter as the paramter to identify the record.

- Expected Content Type: 

```http
    Content-Type: application/json
```

- Expected Request Parameters

```json
    {
        :public_id
    }
```

- Expected Request Body

```json
```

- Expected Status Code

```
Failure: 500
Success: 200
```

- Expected Response Body

```json
{


    "public_id":	"",

    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0

}
```

## GET /Ignitor/

Get retrives all records of the Ignitor type from the HTTP API.

- Expected Content Type: 

```http
    Content-Type: application/json
```

- Expected Request Parameters

```json
    {
    }
```

- Expected Request Body

```json
```

- Expected Status Code

```
Failure: 500
Success: 200
```

- Expected Response Body

```json
[{


    "public_id":	"",

    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0

}]
```

## PUT /Ignitor/:public_id

Update attempts to update a giving record of the Ignitor type from the HTTP API returning received result as a JSON
response. It uses the provided `:public_id` parameter as the paramter to identify the record with the provided JSON request body.

- Expected Content Type: 

```http
    Content-Type: application/json
```

- Expected Request Parameters

```json
    {
        :public_id
    }
```

- Expected Request Body

```json
{


    "name":	"",

    "rex":	{
	
	
	    "url":	""
	
	},

    "rack":	0,

    "public_id":	""

}
```

- Expected Status Code

```
Failure: 500
Success: 204
```


- Expected Response Body

```json
```

## DELETE /Ignitor/:public_id

Get deletes a giving record of the Ignitor type from the HTTP API returning received result as a JSON
response. It uses the provided `:public_id` parameter as the paramter to identify the record.

- Expected Content Type: 

```http
    Content-Type: application/json
```

- Expected Request Parameters

```json
    {
        :public_id
    }
```

- Expected Request Body

```json
```

- Expected Status Code

```
Failure: 500
Success: 204
```

- Expected Response Body

```json
```