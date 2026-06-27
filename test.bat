@echo off
echo Running DBDB tests...
go test -v -cover .\internal\...
if %errorlevel% neq 0 (
    echo Tests failed!
    exit /b %errorlevel%
)
echo All tests passed successfully!
