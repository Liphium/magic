# Magic File Sharing Example

This example is a simple example showing you how to use the SeaweedFS driver to set up a little file sharing backend using the Go standard library and the official AWS S3 SDK Golang SDK.

The project structure here is not how your's should look in actual software that's not just for example purposes. Please look at the `real-project` example for a proper structure (it's one folder up).

## Commands

As with all Magic examples, just use the default `go` CLI to start the example:

```sh
go run .
```

If you run the test we have, you can do (I recommend `-v` because SeaweedFS takes long to start up):

```sh
go test ./... -v
```

## The app in this example

This app just contains a really small HTTP file sharing service that uses SeaweedFS as the storage backend.

- `POST :path` - Upload a file
- `GET :path` - Download a file

If you want to try out this service, I would recommend using the scripts in this repository. Find more info below.

## Scripts

There are various [scripts](https://liphium.dev/magic/documentation/magic-scripts/) available in this example. You can run the app and then upload, download or clear all the files using just scripts. Feel free to try it out. For usage of scripts, please refer to the [documentation](https://liphium.dev/magic/documentation/magic-scripts/).

## Project structure

And finally, to give you another look over what's in this example, here is a quick explanation of all of the files:

```
file-sharing/
├── main.go              # Entrypoint for the program
├── app/
│  ├── app.go            # Main app where the HTTP server is
│  ├── app_test.go       # Simple test for the app's functionality
│  ├── config.go         # Config for Magic
│  ├── scripts_files.go  # Scripts for uploading / downloading files
│  └── storage.go        # Connection to SeaweedFS using the S3 SDK
└── main.go              # Entrypoint for the program
```
