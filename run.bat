@echo off
REM Development quick start script

echo Starting BatchInvoice PDF...

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo Go is not installed. Please install Go 1.21 or later.
    pause
    exit /b 1
)

REM Install dependencies
echo Installing dependencies...
go mod download

REM Run the application
echo Running application...
go run main.go

pause
