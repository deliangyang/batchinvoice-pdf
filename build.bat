@echo off
REM BatchInvoice PDF - Windows Build Script

echo Building BatchInvoice PDF for Windows...

REM Create build directory
if not exist "build" mkdir build

REM Get version info
set VERSION=1.0.0
set BUILD_TIME=%date% %time%
set GIT_COMMIT=unknown

REM Build
echo Building...
go build -ldflags "-s -w -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o build\batchinvoice-pdf-windows.exe main.go

if %errorlevel% == 0 (
    echo.
    echo Build completed successfully!
    echo Output: build\batchinvoice-pdf-windows.exe
    echo.
    echo To run the application:
    echo   build\batchinvoice-pdf-windows.exe
) else (
    echo.
    echo Build failed!
)

pause
