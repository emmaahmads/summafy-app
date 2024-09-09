#!/bin/bash

# Set the URL and headers
URL="http://localhost:8080/notification"
HEADER="Content-Type: application/json"

# Define the data to be sent
DATA=(
  '{"filename": "file1.txt", "summary": "This is file 1"}'
    '{"filename": "file2.txt", "summary": "This is file 2"}'
      '{"filename": "file3.txt", "summary": "This is file 3"}'
        # Add more data here...
)

# Send the requests
for payload in "${DATA[@]}"; do
	  curl -X POST \
		      $URL \
		          -H "$HEADER" \
			      -d "$payload"
		      done
