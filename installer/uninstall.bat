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

REM Start uninstalling
REM Delete only known files. Preserve unrelated files in this directory.
DEL /F "%ProgramW6432%\RDPRelativeInput\RelativeInputClient.exe"
DEL /F "%ProgramW6432%\RDPRelativeInput\RelativeInput.dll"
RMDIR "%ProgramW6432%\RDPRelativeInput"

"%SystemRoot%\System32\reg.exe" DELETE "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" /f

ECHO done.
PAUSE
