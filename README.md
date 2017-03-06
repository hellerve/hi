# hi

A dead simple chat server and client without any external dependencies.

The UI is inspired by material design, but comes without framework bloat.
The entire page is about 19kb in size, without any minification.

The server is written in Go and pretty fast. Also, it's dead-simple.

There are probably bugs.

## Channels

Channels are weird. Anyone can send messages to any channel. A channel
exists as long as there is someone in that channel. If there isn't, the
channel will be closed. What "being in a channel" means is basically
that you subscribe to it, i.e. you will receive messages posted in that
channel. It doesn't affect the way you post there, though. Channels
are conversations. You can scream something at a group of people, or
choose to join that group and interact with them.

## Security

Nothing is stored in a database. Everything is kept in memory, there
is no history. That means that all knowledge and all communication is
ephemeral, like in a party where you join and leave conversations as
you please. Please use HTTPS when deploying this chat, preferably through
nginx or any other web server.

## Deployment

This is where it gets tricky. If you know how to (cross-)compile Go programs
and put them on your server and run it, deploying `hi` is as simple as compiling,
e.g. with `GOOS=<os, probably linux> go build` and copying the resulting binary
and the public directory to a server, then letting it run. Letting it run could
be more or less involved, depending on your hosting provider. I'll be glad to help.

## Running

The resulting Go program can be run without arguments (runs on port 8080) or with
the port set (using the `-p` option). It assumes that the public directory is in
the directory where the program is started.

## Commands

There are a few special commands that you can issue to interact with the server,
IRC-style.

```
/join <channelname>  # join a channel
/leave <channelname> # leave a channel (will send back an error if user is not subscribed to the channel)
/list                # list all users in the current channel
/channels            # list all channels
```


<hr/>

Have fun!
