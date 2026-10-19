# RDPRelativeInput

## About This Program

This program is designed to allow relative input in an RDP session by wrapping an existing remote desktop client window with another window and sending the client's input information using a RDP Virtual Channel. Currently, only sessions from a Windows machine to a Windows machine is supported.

![sample](https://gyazo.com/be1c9e2af08539d06cebe4932b4e568d.gif)

## Installation

### Windows

1. Download server and client programs from Releases<br />
   [Releases](https://github.com/TKMAX777/RDPRelativeInput/releases)
2. Place `install.bat`, `RelativeInputClient.exe`, and `RelativeInput.dll` in the
   same folder on the client machine. Right-click `install.bat` and choose
   **Run as administrator**. 

## Usage

### Connect to Windows

1. Open Remote Desktop Connection and connect to your server like usual 
2. Run `RelativeInputServer.exe` on host
3. Enjoy!

`RelativeInputClient.exe` is launched automatically by the RDP add-in; do not
start it manually. Run `RelativeInputServer.exe` **inside the remote PC's RDP
session**, not on the local client PC. After installing or replacing the add-in,
close Remote Desktop Connection (`mstsc.exe`) and start a new connection so it
loads the installed DLL.

To uninstall, right-click `uninstall.bat` and choose **Run as administrator**.
It removes only the installed EXE and DLL; unrelated directory contents are
preserved.

  ☆ If you need client cursor, use the F8 key to switch to absolute input.<br />
  ☆ To return to relative input mode, select the RDP Input Wrapper window and hit the F8 key again.<br />
  ☆ If you are using a keyboard setting other than the US keyboard setting on the client machine, the response speed may be significantly reduced due to the IME.<br />
      In this case, please add the US keyboard from the Windows settings.<br />
      `Settings -> Time and Language -> Add Language -> English (US) -> Language Options -> Add Keyboard`<br />
  ☆ When you want to quit this program, please press and hold F12.

## Windows mouse input

Relative mouse motion is captured with `RegisterRawInputDevices` and
`WM_INPUT` / `GetRawInputData`. The client forwards device deltas directly,
without reading and recentering the OS cursor for every movement. Legacy
`WM_MOUSEMOVE` is ignored to prevent duplicate movement; keyboard, button, and
wheel messages retain their existing handling. Protocol v3 remains unchanged.

Raw input registration is foreground-only. Relative mode confines the cursor
to the wrapper's center; switching mode or losing focus releases confinement.
Absolute raw devices establish a baseline before sending position differences.
Client-side Windows pointer acceleration no longer affects relative motion,
so mouse sensitivity may feel different. The server continues to use
`SendInput`; network and application latency are not removed by this change.

## Build

MinGW installation is required to create add-ins for mstsc.exe.

Clone this repository and run these commands on powershell.

```powershell
go build -ldflags -H=windowsgui .\cmd\RelativeInputClient
go build -ldflags -H=windowsgui .\cmd\RelativeInputServer
go build -buildmode=c-shared  -o .\RelativeInput.dll .\windows\virtualchanneladdin
Copy-Item installer\install.bat .
Copy-Item installer\uninstall.bat .
```

After building, follow the [installation instructions](#installation) to install
the client EXE and DLL.

## Server-side cursor overlay

The server draws a software cursor in a transparent, click-through window. RDP
can transmit it as desktop pixels without relying on the Graphics Capture API
cursor-visibility side effect. Replace `RelativeInputServer.exe` on the remote
computer and restart it; this feature requires no protocol or client update.

Click the English cursor button to cycle through these modes:

- **Cursor: Auto** (default): follow the visibility reported by Windows.
- **Cursor: Always**: draw the available cursor even when Windows reports it as
  hidden. Try this if Auto shows no cursor over RDP. This also overrides intentional
  hiding by applications such as games.
- **Cursor: Off**: disable the overlay, for example when another cursor is visible.

The overlay uses physical screen coordinates, including cursor hotspots and
monitors with negative coordinates. It does not take focus or receive clicks.
Updates use a 16 ms UI timer without changing the system timer resolution; actual
latency depends on Windows scheduling and RDP screen updates.

This is a color-key overlay: magenta pixels are transparent, translucent edges
and monochrome invert cursors may differ from native rendering, and animated
cursors use their first frame. Secure desktops, exclusive fullscreen applications,
and RDP paths that exclude layered windows may not display the overlay. Always
mode requires Windows to supply a valid cursor shape. If the native RDP cursor
is also visible, select Off to avoid two cursors.

## Known issues

- Windows Defender misidentifies the DLL as a Trojan horse program. (https://github.com/TKMAX777/RDPRelativeInput/issues/2)

## Copyright

Copyright 2022- tkmax777 and contributors
