# Annotations
This package provides a list of annotations which provide code generation for specific needs and functionality required for a type or a series of types. Each supported annotation is meant to provide a need and also to allow the generation of pices or whole packages, which allow ease of development and also provide a consistent set of API principles coded into this generators for each annotations, that allow a plug-and-play style approach in the areas they are used.

*We hope to expand these list of annotations has time passes and also create a consistent means of easily code-generating parts or whole API.*

## All (Supported/To-Be Supported) Annotations

### @assets

*@assets* provides a package-level annotation for generating a sub-package which will allow you to quickly generate a internal go file within that sub-package to contain all files associated with a given set of extensions as part of a go file. It provides a tear-down `bind-data` where all sources are copied as part of a registery for access.

Example can be found in `Dap` package: [Dap](../examples/dap)

### @templator

*@templator* provides a template-base annotation for generating new code from the provided templates associated with the annotation and from associated annotation `templaterTypesFor` which details necessary key-value pairs for the code generation to be done.

Example can be found in `Temples` package: [Temples](../examples/temples)

 ### @httpapi

*@httpapi* provides a struct-level annotation for generating a http based CRUD API, which provides all necessary calls for the creation, update, removal and retreival of that given struct. This provides a means to quickly generate out a very quick API with basic API readme, that describes each endpoint.

Example can be found in `Dap` package: [Dap](../examples/dap)

### @mongo

*@mongo* provides a struct-level annotation for generating a mongo based type which exposes a `.Exec` method to execute transactions against a mongo collection.

Example can be found in `Dap` package: [Dap](../examples/dap)

### @iface

*@iface* provides a interface-level annotation for generating initial struct implementation with mocking and initial tests files. It's great to allow one
quickly reduce repetitive tasks of writing method implementations for an interface.


### @mongoapi

*@mongoapi* provides a struct-level annotation for generating a mongo based CRUD API for saving a giving type of Struct into a mongo database. This provides a simple and elegant approach where the details of a struct declaration and it's tags defined the name for the record within the record to be saved.

The `@mongoapi` annotation provides an associative annotation which helps customize specific areas of a it's code generation:

    - Create Struct Type (`@associates(@mongoapi, New, StructType)`)
    - Update Struct Type (`@associates(@mongoapi, Update, StructType)`)

    where:
        `@associates` dictates the following relation:
         1. It is to a `@mongoapi` annotation
         2. Its Action is a `New` or `Update`
         3. To use the `StructType` name represented as the type to be expected by either C or U in CRUD.

This two types above are declared with an `@associates` annotation, which will dictate the struct to be presented to either the `Create` and `Update` methods as the struct type to be giving to the function to retrieve the update values from. Generally if not declared the default struct type which the `@mongoapi` annotation is declared on.

Example can be found in `Dap` package: [Dap](../examples/dap)

### @sqlapi

*@sqlapi* provides a struct-level annotation for generating a SQL based CRUD API for saving a giving type of Struct into a SQL database. This provides a simple and elegant approach where the details of a struct declaration and it's tags defined the name for the record within the record to be saved.

The `@sqlapi` annotation provides an associative annotation which helps customize specific areas of a it's code generation:

    - Create Struct Type (`@associates(@sqlapi, New, StructType)`)
    - Update Struct Type (`@associates(@sqlapi, Update, StructType)`)

    where:
        `@associates` dictates the following relation:
         1. It is to a `@sqlapi` annotation
         2. Its Action is a `New` or `Update`
         3. To use the `StructType` name represented as the type to be expected by either C or U in CRUD.

This two types above are declared with an `@associates` annotation, which will dictate the struct to be presented to either the `Create` and `Update` methods as the struct type to be giving to the function to retrieve the update values from. Generally if not declared the default struct type which the `@sqlapi` annotation is declared on.


Example can be found in `Dap` package: [Dap](../examples/dap)
