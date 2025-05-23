# Liphium Magic: Test your applications with confidence.

This project contains lots of experimental tools for database testing and more. None of the tools in this repository are fully featured and tested, please use with caution and do not use in mission-critical projects.

The goal of Liphium Magic is to built a testing toolkit so powerful that testing software is actually fun. Unit testing is easy and can be a nice way to test your projects. I want to make testing your complex backend just as easy as unit testing.

## Hackathon Goals

Folder structure:
- .magic
    - .gitignore (contains database folder and local persistence)
    - config.go (synced, for configuration what the app needs)
    - local.json (contains default profile for running the app with ``magic debug``)
    - databases/ (for all databases, locally persisted)
    - scripts/ (all scripts, synced with repository)
    - tests/

### Definitely

- General CLI
    - Help command
    - ``magic init`` for generating gitignore and basic configuration
- Magic SDK
    - Control databases and stuff (and create connections)
    - Communicate with the Runner
- Magic Config SDK
    - Can run the Runner
- Magic Runner
    - Run the project, database containers, etc.
- Debugging of local apps
    - Create databases from configuration using Docker
    - Create local persistence in the local.json file (for databases used in the app)
    - Database management using the magic CLI

### If there is time

- Docker configuration and setup
    - Build a Docker image from the current repository (using Dockerfile)
- Base for testing and scripting
    - A directory with files/directories for tests and scripts written in Go
    - Do code generation that can execute the functions in those tests


Wichtig: https://excalidraw.com/#json=RhQtkM0Ufq11Sf7CrvYMa,YdBfJ9cxWAhoMhTQ1AYgNw