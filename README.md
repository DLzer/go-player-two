# Go Player Two

Go Player Two is a small yet robust multiplayer game server written in go and meant for use in light bi-directional websocket communication between two parties.

# Building

## Build using Makefile

Make a new server
```bash
make aserver
>> Starting Server...
```

Make a new client to test the server. Two clients makes a party.
```bash
>> Starting Client...
```

## Build the binary

Select the binary that matches your OS from the latest release version.
Run the binary to start the server. By default connections are over port 40000.

# Operations

Starting the server will spin up an http server and listen for connections hitting the `/engine` endpoint. A matchmaking pool will already be awaiting client connections. Once the matchmaking pool has two client it will start the bi-directional communication of 'game' data.

Clients can send data through the upgraded connection to be verified and passed through the game server which will update the opponents player state.

# Disclaimer

This project is a work in progress and is not meant for marketable or production use.
