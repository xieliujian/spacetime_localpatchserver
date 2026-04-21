@echo off
echo Building patch server...

set GOOS=windows
set GOARCH=amd64

"D:\Program Files\Go\bin\go.exe" build -o patchserver.exe cmd/server/main.go

if %errorlevel% == 0 (
    echo Build successful: patchserver.exe
) else (
    echo Build failed
    pause
)
