# RDPRelativeInput

## About This Program

This program is designed to allow relative input in an RDP session by wrapping an existing remote desktop client window with another window and sending the client's input information using a RDP Virtual Channel. Currently, only sessions from a Windows machine to a Windows machine is supported.

![sample](https://gyazo.com/be1c9e2af08539d06cebe4932b4e568d.gif)

## install

### Windows

1. Download server and client programs from Releases`<br>`
   [Releases](https://github.com/TKMAX777/RDPRelativeInput/releases)
2. Run `install.bat` on the client machine.

## Usage

### Connect to Windows

1. Open Remote Desktop Connection and connect to your server like usual and have it in Maximize Windowed (**NOT FULL SCREEN**)
2. Run `RelativeInputServer.exe` on host
3. Enjoy!

  ☆ The mouse cursor disappears during relative input mode. If you need the cursor, use the F8 key to switch to absolute input.`<br />`
  ☆ To return to relative input mode, select the RDP Input Wrapper window and hit the F8 key again.`<br />`
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

## Copyright

Copyright 2022- tkmax777 and contributors
