@echo off
cd /d "%~dp0"
go test ./repository -v > test_results.txt 2>&1
type test_results.txt
