# Qubes Updater (Go)

A simple, effective Go-based CLI tool to manage and check for updates across your Qubes OS virtual machines. It ensures your templates and VMs stay updated — with support for update exclusions and optional auto-update.

## Features

- Lists all available Qubes VMs (excluding `dom0`)
- Checks for available updates in each VM using `dnf`
- Skips updates for excluded VMs (defined in a file)
- Optionally auto-updates all VMs with pending updates
- Logs all actions and errors to `qubes-updater.log`

## Requirements

- Must be run from `dom0` in Qubes OS  
  Qubes commands like `qvm-run` and `qvm-ls` are only available in dom0.
- Requires Go 1.20+ to compile

## Installation

1. Clone the repo or copy the code
2. Create a Go module (if needed):

```
go mod init qubes-updater
```

3. Build it:

```
go build -o qubes-updater main.go
```

## Usage

Check for updates only:

```
./qubes-updater
```

Check and automatically update:

```
./qubes-updater --update
```

Use a custom exclusion file:

```
./qubes-updater --exclude=my-vm-exclusions.txt
```

## Excluding VMs from updates

To skip updates for specific VMs (e.g. `firefox`, `vault`, `vpn`), create a file called `exclude.txt` in the same directory:

```
firefox
vault
vpn
```

You can specify a different file using the `--exclude` flag.

## Log file

All actions and errors are written to `qubes-updater.log`.  
This helps you track:

- Which VMs had updates
- Which were skipped
- Any errors encountered
- Output from update attempts

## Important Notes

- Run from dom0. The script will not work inside other VMs (`work`, `personal`, etc).
- So u need to do :
  ```` qvm-copy-to-vm dom0 qubes-updater ````
- It uses `sudo dnf updateinfo list available` for checking updates.
- You can customize it to use `apt` if you're running Debian-based VMs.

## Planned Improvements

- Interactive terminal UI (TUI)
- Per-VM update options
- Real-time update progress

## Contributing

Pull requests and feedback welcome.  
Feel free to fork and adapt the script to your environment.

## License

MIT – Use freely and improve it as you wish.
