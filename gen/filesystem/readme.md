FileSystem
-----------------
FileSystem provides a simple API to create an in-memory file system with internal
directories and files, represented by a composition of functions.


## API
To create an in-memory file system with simple

```go
dockerFS = FileSystem(
  Meta("version", "1.0"),
  Description("FileSystem for the faas docker system"),
  Dir("app"),
  File(

  ),
  File(
    "dockerfile",
    Content(`
      FROM alpine:latest

      RUN sudo apt-get update
      RUN sudo apt-get install git golang

      CMD ["/bin/bomb"]
    `)
  ),
)

var dest io.Writer

gzipDockerFS = GZipTarFS(dockerFS)
gzipDockerFS.WriteTo(dest)

tarDockerFS = TarFS(dockerFS)
tarDockerFS.WriteTo(dest)

zipDockerFS = ZipFS(dockerFS)
zipDockerFS.WriteTo(dest)
```
