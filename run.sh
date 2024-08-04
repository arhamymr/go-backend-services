
#!/bin/bash

# Check if Redis server is running
if ! pgrep -x "redis-server" > /dev/null
then
    echo "Starting Redis server..."
    redis-server &
else
    echo "Redis server is already running."
fi


nodemon --exec go run cmd/main.go --signal SIGTERM