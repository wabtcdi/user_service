#!/bin/bash
cd /c/Users/19802/GolandProjects/user_service
echo "Running service tests..."
go test -v ./service 2>&1
echo "Tests completed with exit code: $?"
