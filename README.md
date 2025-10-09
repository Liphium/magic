# Liphium Magic: Database testing for everyone.

This project contains lots of experimental tools for database testing and more. None of the tools in this repository are fully featured and tested, please use with caution and do not use in mission-critical projects.

The goal of Liphium Magic is to built tools to make testing applications written in Go that rely on a database connection to work easy. Unit testing was always easy and can be a nice way to test your projects if you have just one function to test. Magic makes testing your complex logic and backend as easy as unit testing by leveraging Docker and code generation.

## Usage

**1.** Install Magic using the following command:

```sh
go install github.com/Liphium/magic/v2@latest
```

**2.** To use Magic go to any Go project of your chosing and run `magic init`. This will create a new directory in your project called `magic`. In there you can edit the `config.go` file to create databases. Here's an example:

```go
// ...

// This is the function called once you run the project
func Run(ctx *mconfig.Context) {
	// Create a new PostgreSQL database called main
	db := mconfig.NewPostgresDatabase("main")
	ctx.AddDatabase(db)

	// Add environment variables so you can access the database later
	ctx.WithEnvironment(&mconfig.Environment{
		"DB_HOST":     db.Host(ctx),
		"DB_PORT":     db.Port(ctx),
		"DB_USERNAME": db.Username(),
		"DB_PASSWORD": db.Password(),
		"DB_NAME":     db.DatabaseName(ctx),
	})
}

func Start() {
    // Start your application here (you may have to rename your main function or move it to a different module, sorry, otherwise Magic can't work)
}

// ...
```

**3.** You can now use `magic start` to run your app. Become a great wizzard!

## Status

Magic is still in very early development. We'll continue developing it until it can be used to test our own backend which is going to still require some complex problems to be solved. However, we already see a lot of potential in this tool to provide a really nice developer experience.

Current version: v1.0.0-rc12.
