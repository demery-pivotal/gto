`gto` is is a very crude tool for separating a Gradle `output.bin` file into
stdout and stderr files for each test method. `gto` stands for Gradle Test
Output (and also a Pontiac muscle car that was popular in a previous
millenneum).

`gto` parses a file called `output.bin` in the current working directory.

It creates a directory called `results` and writes:
- a `results/class-CCCCC/method-MMMMM/stdout.log` file for every test method
  that writes to stdout
- a `results/class-CCCCC/method-MMMMM/stderr.log` file for every test method
  that writes to stderr

# Installation
`cd` to the repo's root directory and run `go install .` This will install a
`gto` executable in your go installation's `bin` directory (typically
`~/go/bin`.

# Limitations
- It always ends with an `EOF` error, even though it ought to expect
  `output.bin` to have a finite length.
- If a target output file already exists, such as from a previous execution,
  `gto` will append to it rather than overwriting it.
- `gto` currently offers no CLI options. It would be nicer if it offered
  options to specify the input file and the output directory.
- `gto` reads only the `output.bin` file, which knows classes and methods only
  by numeric IDs. It would be nicer if it detected the `output.bin.idx` file
  (if it exists) and used that to name its class and method directories.
