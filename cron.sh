#!/usr/bin/bash
export CURRENT_FILE=$(realpath $0)
export CURRENT_PATH=$(dirname $CURRENT_FILE)
cd "$CURRENT_PATH"
go run main.go
git push