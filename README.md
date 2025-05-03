# Liphium Magic: Test your applications with confidence.

This project contains lots of experimental tools for database testing and more. None of the tools in this repository are fully featured and tested, please use with caution and do not use in mission-critical projects.

The goal of Liphium Magic is to built a testing toolkit so powerful that testing software is actually fun. Unit testing is easy and can be a nice way to test your projects. I want to make testing your complex backend just as easy as unit testing.

### Idea

Magic is a test runner.

It will take in a config (.magic/config.yml) and build a Docker setup completely automatically.

To add a test you can create a directory like this: .magic/tests/test_name or just create a .magic/tests/test_name.yml file. If you want to include stuff like a .sql file or something, use the directory. If you just need a test script, use the test_name.yml. For the directory the test script will be .magic/tests/test_name/test.yml.

The test currently being executed will be in the environment variable MAGIC_TEST. Through that your regular test code can identify what needs to be run. If the environment variable isn't the test name, the test should not be included for the run as regular unit tests should just run normally. Only the test code should run when the environment variable is set to prevent duplicate runs.

We should also make it so you can specify other environment variables just to make things simpler. Say you want to run the complete project in two tests and not in one other test. It would be annoying having to always switch case or use if statements for such simple things. So for less code repetition we could add that you can set your own environment variables as well. Maybe only allow prefixed with MAGIC\_?

Things like wait points and parallel execution are gonna be important for the test runs... We could add something like waitpoints and then have a waitgroup under the hood. Where for example in 1 test you could have two parallel executions: One then says wait for "waitpoint1" to become 2 while the other is still executing something and then also waits for "waitpoint1" to become 2. Because it's been called twice it's now two and both start executing at the same time again. Things like just waiting for a little bit could also be good. Or pinging and endpoint until it is available or something and then continuing the test run.
