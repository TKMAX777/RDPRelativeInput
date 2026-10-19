@echo off
setlocal
REM Require explicit administrator execution; never execute a temporary script.
"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -NonInteractive -Command "if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) { exit 1 }"
if errorlevel 1 (
    echo Right-click this file and select Run as administrator.
    pause
    exit /b 1
)
if not defined ProgramW6432 exit /b 1
:gotAdmin
    pushd "%CD%"
    CD /D "%~dp0"

if not exist "RelativeInputClient.exe" exit /b 1
if not exist "RelativeInput.dll" exit /b 1
REM Start installing

REM Create directory
MKDIR "%ProgramW6432%\RDPRelativeInput"

REM Install it
COPY /Y RelativeInputClient.exe "%ProgramW6432%\RDPRelativeInput\RelativeInputClient.exe"
if errorlevel 1 exit /b 1
COPY /Y RelativeInput.dll "%ProgramW6432%\RDPRelativeInput\RelativeInput.dll"
if errorlevel 1 exit /b 1

@REM REG ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" 
"%SystemRoot%\System32\reg.exe" ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" /f
"%SystemRoot%\System32\reg.exe" ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" /v Name /t REG_SZ /f /d "%ProgramW6432%\RDPRelativeInput\RelativeInput.dll"
if errorlevel 1 exit /b 1

ECHO done.
PAUSE
