@echo off
cd %~dp0
start "ONE-AIR" oneair_server.exe
start http://127.0.0.1:8080
