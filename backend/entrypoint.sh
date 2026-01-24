#!/bin/sh

# Start Redis in the background
redis-server --daemonize yes

# Start the Go application
./main
