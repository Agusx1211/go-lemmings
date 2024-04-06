# Go-Lemmings

Go-Lemmings is a simple, yet effective tool designed to manage process spawning based on the current system load. It observes the system load and spawns processes from a specified template, maintaining the load around a user-defined target. The tool is useful for system load testing or managing batch jobs without overloading the system.

## Building and Installing

Ensure you have a Go environment set up, then follow these steps:

1. Clone the source code (assuming you have it in a Git repository):
   ```bash
   git clone https://your-repository-url/go-lemmings.git
   cd go-lemmings
   ```

2. Compile the program:
   ```bash
   go build -o go-lemmings main.go
   ```

3. Optionally, install the program to your system:
   ```bash
   go install
   ```

## Usage

To use Go-Lemmings, execute the following command with the appropriate flags:

```
./go-lemmings -l [target load] -m [max processes] [command template]
```

- `-l [target load]`: Desired system load level. Go-Lemmings will spawn processes to try and reach this load.
- `-m [max processes]`: The maximum number of processes that can be spawned.
- `[command template]`: The command to be executed as new processes, where `{random}` can be used to inject a random number.

### Example

To keep the system load around 5 with no more than 250 processes executing a sleep command, run:

```bash
./go-lemmings -l 5 -m 250 sleep 0.1
```

For commands that should include a random component, use `{random}` in the template:

```bash
./go-lemmings -l 5 -m 250 'echo {random} && sleep 0.1'
```

Go-Lemmings will replace `{random}` with a random uint32 integer each time it initiates a process.

### Notes

- Ensure that both the target load and maximum processes are greater than zero.
- Enclose the command template in quotes if it contains spaces or special characters to ensure it is passed correctly to the program.