# RDPRelativeInput

## About This Program

This program is designed to allow relative input in an RDP session by wrapping an existing remote desktop client window with another window and sending the client's input information using a RDP Virtual Channel. Currently, only sessions from a Windows machine to a Windows machine is supported.

![sample](https://gyazo.com/be1c9e2af08539d06cebe4932b4e568d.gif)

## install

### Windows

1. Download server and client programs from Releases<br />
   [Releases](https://github.com/TKMAX777/RDPRelativeInput/releases)
2. Run `install.bat` on the client machine.

## Usage

### Connect to Windows

1. Open Remote Desktop Connection and connect to your server like usual 
2. Run `RelativeInputServer.exe` on host
3. Enjoy!

  ☆ If you need client cursor, use the F8 key to switch to absolute input.<br />
  ☆ To return to relative input mode, select the RDP Input Wrapper window and hit the F8 key again.<br />
  ☆ To close this app, hold the F12 key.<br />
  ☆ Administrator privileges are required for operation in some games. In that case, please run RelativeInputServer.exe with Administrator privileges.

## Build

MinGW installation is required to create add-ins for mstsc.exe.

Clone this repository and run these commands on powershell.

```powershell
go build -ldflags -H=windowsgui .\cmd\RelativeInputClient
go build -ldflags -H=windowsgui .\cmd\RelativeInputServer
go build -buildmode=c-shared  -o .\RelativeInput.dll .\windows\virtualchanneladdin
cp installer\install.bat
.\install.bat
```

## Why this program uses Graphics Capture API?

This program uses Graphics Capture API on the host computer side.

In a remote desktop connection, the cursor on the host side is normally not displayed. However, the GraphicsCaptureAPI allows the cursor to be displayed, so this program creates a button which can turn this feature on and off and captures that window.

## Copyright

Copyright 2022- tkmax777 and contributors
