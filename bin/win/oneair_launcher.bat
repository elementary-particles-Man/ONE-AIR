@echo off
cd /d "%~dp0"
echo Starting ONE-AIR...
start "" oneair_server.exe
timeout /t 2 >nul
start "" http://127.0.0.1:8800/
