@echo off

REM Check if this session has administrative privileges.

>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"

if '%errorlevel%' NEQ '0' (
    echo Requesting administrative privileges...
    goto UACPrompt
) else ( goto gotAdmin )

:UACPrompt
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
    set params = %*:"="
    echo UAC.ShellExecute "cmd.exe", "/c %~s0 %params%", "", "runas", 1 >> "%temp%\getadmin.vbs"

    "%temp%\getadmin.vbs"
    del "%temp%\getadmin.vbs"
    exit /B

:gotAdmin
    pushd "%CD%"
    CD /D "%~dp0"

REM Start installing

REM Create directory
MKDIR "%ProgramW6432%\RDPRelativeInput"

REM Install it
COPY RelativeInputClient.exe "%ProgramW6432%\RDPRelativeInput\RelativeInputClient.exe"
COPY RelativeInput.dll "%ProgramW6432%\RDPRelativeInput\RelativeInput.dll"

@REM REG ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" 
REG ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" /f
REG ADD "HKCU\Software\Microsoft\Terminal Server Client\Default\AddIns\RelativeInput" /v Name /t REG_SZ /f /d "%ProgramW6432%\RDPRelativeInput\RelativeInput.dll"

ECHO done.
PAUSE
