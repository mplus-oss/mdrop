# mdrop

mdrop is a command-line tool that facilitates secure file transfer between two computers using a peer-to-peer (P2P) tunneling mechanism established over an SSH connection. This approach draws inspiration from popular solutions like ShareDrop and LocalDrop.

## Download

Currently, mdrop is available for download on the GitHub Releases page. Pre-built binaries are currently only supported for Linux (x86_64 architecture).

## Installation

The installation process for mdrop depends on the chosen method of brokering the connection: Client-side or Docker image based.

### Prerequisite (Client)
1. SSH Client
2. Shell (sh)
3. Optional Proxy Support:
   - Cloudflare WARP is currently supported for proxy connections (cloudflared).

### Prerequisite (Manual Broker)
1. SSHD Daemon

### Prerequisite (Docker Image Broker)
To leverage a Docker container as the broker, pull the following image:
```sh
docker pull ghcr.io/mplus-oss/mdrop-sshd-tunnel:latest
```

## Usage

### Authentication
This command initiates the authentication process for secure file transfer. Follow the on-screen prompts to set up user credentials.
```
mdrop auth
```

### Sending Files
Use this command to transfer files to another computer. Options can be accessed with `mdrop send --help`.
```
mdrop send [options] <file1> [file2] [file3] ...
```

### Receiving Files
This command facilitates receiving files from another computer. Explore available options with `mdrop get --help`.
```
mdrop get [options] <token>
```
