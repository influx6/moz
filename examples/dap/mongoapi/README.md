Ignitor DB API 
===============================

Ignitor DB API attempts to provide a simple basic documentation which details 
the basic structure of the Ignitor type, and the response which will be provided 
when working with this API.

The API provides a basic CRUD interface as describe below:

## Create

Create stores a given Ignitor type into the mongo db, where the DB API expects the provided type to be called
against the API function type for adding a given record Ignitor.

```go
Create(ctx context.Context, elem dap.Ignitor) error 
```

## Get

Get retrives a giving record of the Ignitor type from the DB API returning received result matching 
the publicID value provided if found else returning an error.

```go
Get(ctx context.Context, publicID string) (dap.Ignitor, error) 
```

## Get All

Get retrives all records of the Ignitor type from the DB API.

```go
GetAll(ctx context.Context) ([]dap.Ignitor, error) 
```

## Update

Update stores a given Ignitor type into the mongo db, where the DB API expects the provided type to be called
against the API function type for updating a existing record Ignitor.

```go
Update(ctx context.Context, elem dap.Ignitor) error 
```

## Delete

Delete destroys the giving record of the Ignitor type from the DB API returning an error if any occured.

```
Delete(ctx context.Context, publicID string) error 
```