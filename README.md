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
Run Go-Lemmings with the following command:

```
./go-lemmings -l [target load] -m [max processes] [-d delay ms] [-i ignore err] [command template]
```

Options:
- `-l [target load]`: The target system load that Go-Lemmings will attempt to maintain.
- `-m [max processes]`: The maximum number of processes that Go-Lemmings will spawn.
- `-d [delay ms]`: (Optional) Delay in milliseconds to wait after starting each process.
- `-i [ignore err]`: (Optional) A substring of the error output to ignore. If this substring is found in the error output of a process, that error will be logged but not cause all processes to terminate.

`[command template]`: The command that Go-Lemmings will run as new processes. Use `{random}` in the template to insert a random uint32 value into each command instance.

### Examples

- To maintain a system load of 5 with a maximum of 250 processes, each executing a command that sleeps for 0.1 seconds:
  ```bash
  ./go-lemmings -l 5 -m 250 sleep 0.1
  ```

- To run the same setup with a delay of 100 milliseconds between process starts:
  ```bash
  ./go-lemmings -l 5 -m 250 -d 100 sleep 0.1
  ```

- To ignore errors containing the substring "timeout" while running the commands:
  ```bash
  ./go-lemmings -l 5 -m 250 -i timeout 'ping -c 1 google.com'
  ```

Errors that contain the specified string in `-i` will be logged, allowing Go-Lemmings to continue operation despite these errors.

- To keep the system load around 5 with no more than 250 processes executing a sleep command, run:
  ```bash
  ./go-lemmings -l 5 -m 250 sleep 0.1
  ```

- For commands that should include a random component, use `{random}` in the template:
  ```bash
  ./go-lemmings -l 5 -m 250 'echo {random} && sleep 0.1'
  ```

Go-Lemmings will replace `{random}` with a random uint32 integer each time it initiates a process.

### Notes

- Ensure that both the target load and maximum processes are greater than zero.
- Enclose the command template in quotes if it contains spaces or special characters to ensure it is passed correctly to the program.