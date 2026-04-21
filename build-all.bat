@echo off
echo Building patch server for multiple platforms...

echo.
echo [1/3] Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
"D:\Program Files\Go\bin\go.exe" build -o bin/patchserver-windows-amd64.exe cmd/server/main.go
if %errorlevel% == 0 (echo   Success) else (echo   Failed)

echo.
echo [2/3] Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
"D:\Program Files\Go\bin\go.exe" build -o bin/patchserver-linux-amd64 cmd/server/main.go
if %errorlevel% == 0 (echo   Success) else (echo   Failed)

echo.
echo [3/3] Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
"D:\Program Files\Go\bin\go.exe" build -o bin/patchserver-darwin-amd64 cmd/server/main.go
if %errorlevel% == 0 (echo   Success) else (echo   Failed)

echo.
echo All builds completed. Check bin/ directory.
pause
