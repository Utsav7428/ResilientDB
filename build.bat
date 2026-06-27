@echo off
setlocal
echo Building DBDB...
if not exist "bin" mkdir bin
go build -o bin\dbdb.exe .\cmd\dbdb
if %errorlevel% neq 0 (
    echo Build failed!
    exit /b %errorlevel%
)
echo Build successful. Executable located in bin\dbdb.exe.
endlocal
