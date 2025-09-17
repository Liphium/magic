# Magic Real Project Example

This example showcases what a Magic setup in a real project could look like. While you may want to customize it for your own setup, this is how we (from Liphium) generally structure our projects when working with Magic. We hope this provides a clear reference of how to go about taking full advantage of Magic's entire feature-set.

## Starting the example

The magic about Magic is that, well, all you need is Golang and Docker installed. And then everything will just work with this little command:

```bash
go run .
```

Yes, that's the default run command for your Go app. That's because Magic integrates directly into your project. And if you want to build your project without it, you can use build tags:

```bash
go run -tags release .
```

The same thing works for `go build` as well. You can look at `main_magic.go` and `main.go` for how we implemented this or look up what Go's build tags are if you've never seen them before. Which I honestly won't blame you for, since I've also learned about them while trying to work out how to properly structure an app with Magic xd

## Project structure

Now that you've run the project and seen the magic of Magic, maybe you want to explore the example a little better. Everything is nicely commented to guide you a little bit, but here's an overview of all the files in this project anyway:

```
real-project/
├── main.go              # Production entry point (build tag: release)
├── main_magic.go        # Development entry point with Magic integration
├── ...
├── database/
│   └── database.go    # Database connection and Post model
├── starter/
│   ├── config.go              # Magic configuration and setup
│   ├── start.go               # Web server with Fiber endpoints
│   ├── start_test.go          # Integration tests leveraging Magic's test runner
│   ├── scripts_database.go    # Database management scripts
│   └── scripts_endpoints.go   # API testing scripts
└── util/
    └── requests.go     # HTTP utility functions
```

## The app in this example

This app just contains a really small posts service that can create new posts and list all posts or individual posts using a PostgreSQL database.

- `POST /posts` - Create post
- `GET /posts` - List posts
- `GET /posts/:id` - Get post

While you now could go and get out your API client, Magic has something better: Scripts. Read below to see how you can use them.

## Scripts

Before worrying about anything and looking into their code, just try using them and understanding them after. That might honestly be a better way, since it will teach you how easy it can be for new people to go into a codebase they've never seen before. In fact, why don't you let Magic tell you which scripts there are instead of me wasting time explaining them here:

```bash
go run . --scripts
```

You can then run the scripts using:

```bash
go run . -r script [arguments]
```

You can pass arguments in by either typing them in the CLI, which might be easier for new people, or you can just provide them as arguments behind (in order of the fields in the struct). Magic will also tell you if you mess up.

## Testing

Well to the next magical thing Magic can do for you: Testing. Imagine just being able to start your app with one line of code and then being able to just throw requests at it, and even be able to check the database. Wouldn't that be awesome?

Well, with Magic you can just do this:

```go
func TestMain(m *testing.M) {
    magic.PrepareTesting(m, starter.BuildMagicConfig())
}
```

And you're done. Look into `start_test.go` for how we do it in this little example. Point is, even in this example you just downloaded, you can just do:

```bash
go test ./...
```

And all of the tests will just run by themselves. They can _currently_ not run in parallel, but we're working on a solution for it that should come out in some future update.
