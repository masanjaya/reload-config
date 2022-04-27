#!/bin/bash
set -x
set -e
GOOS=linux go build
scp reload-config tv:/var/tmp/cryptool/v1
