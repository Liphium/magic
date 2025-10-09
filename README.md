# Liphium Magic

## This project is still experimental and in early stages. Feel free to test it out, but expect major changes, bugs and of course also new features.

Liphium Magic is a suite of tools to help build tests and provide a better developer experience when developing web services in Golang. We made it because we felt like it the barrier of making a contribution to our own projects was too high. When working on applications in a team, it's important that everyone can easily start the project and also use the same tools everyone else is using. Like with dependencies, when you need a new database for something, you shouldn't have to tell everyone in your team to complete extra steps just for their app to run again. When someone first joins your project, they can ideally set everything up with one command. That's the vision of Magic, our all-in-one developer experience toolkit. For testing your app, both automatically and manually, as well as running it on your own machine without complex setup.

The path to this goal is of course a long one, and we also know that, so for now Magic can only really help with PostgreSQL databases. It only supports this one simple database and not more. It's all we use in our apps and in the future, when also have the need for it, we will likely integrate more services into Magic or build a nice abstraction layer that makes it easy to integrate different services.

We hope you'll enjoy this project, feel free to open an issue in case anything annoys you or you see potential for a new feature. Be nice and maybe we'll add it. We're working on this in our freetime though, so don't expect us to answer all. Before creating a pull request for something, please consult us in the issues. We're looking forward to having a discussion with you.

## System requirements

- Desktop operating system (Windows, macOS or Linux)
- Docker (must be installed and the Go toolchain must have permissions to access the socket)
- Golang (you're not making a Go application without it)

## Application limitations

Magic only supports specific services, and while we do plan on increasing the amount of supported services, for now we only support the services listed below. If your application needs anything else, you're currently not the target audience for Magic.

### Supported databases

- PostgreSQL

Other services may be supported in the future.

## Features

- Make your app runnable with one command on any machine that meets the System requirements
- Develop scripts that interact with your application or the database
  - Allows sharing of tools you're using for testing
- Test your application using integration tests (they can also call your scripts)
  - Test with a real database using a real connection 

## Usage

**1.** Add Magic to your project:

```sh
go get -u github.com/Liphium/magic/v2@latest
```

**2.** Wrap your main function with ``magic.Start`` (please take a look at the [real project example](https://github.com/Liphium/magic/tree/main/examples/real-project) for how to really to do this, this just serves as a showcase):

```go
// ...

func main() {
	magic.Start(magic.Config{
		AppName: "magic-example",
		PlanDeployment: func(ctx *mconfig.Context) {
			// Create a PostgreSQL database for the posts service
			postsDB := ctx.NewPostgresDatabase("posts")

			// Set up environment variables for the application
			ctx.WithEnvironment(mconfig.Environment{
				// Database connection environment variables
				"DB_HOST":     postsDB.Host(ctx),
				"DB_PORT":     postsDB.Port(ctx),
				"DB_USER":     postsDB.Username(),
				"DB_PASSWORD": postsDB.Password(),
				"DB_DATABASE": postsDB.DatabaseName(ctx),
			})
		},
		StartFunction: Start,
	})
}

func Start() {
    // Start your application here (we have to take over your main function to be able to run code before)
}

// ...
```

**3.** You can now use `go run .` to run your app and a database will be created in a Docker container near you. 

Become a great wizard! If you want to be a real great one though, I would take a look at the [real project example](https://github.com/Liphium/magic/tree/main/examples/real-project) to actually see how it's done.

## Documentation

Because this project is new, there is no real documentation yet. However, you can look into the `examples` folder. It should contain about everything you need for now.
