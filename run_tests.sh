#!/bin/bash
cd "$(dirname "$0")"
go test ./repository/... -v -count=1
