webtalk
=======

The `webtalk` project is a simple application to allow platforms implement
WebSockets easily.

It works as a separate server in your implementation in which you "talk" to
this server to deliver messages to the users. The users must communicate
with your server via your application, as WebTalk only receives the initial
Login message from the users and nothing more.

For the moment it only supports WebSockets because they are clean and easy,
and this was just an experiment that worked really good.

Next steps, maybe implement [go-socket.io](https://github.com/madari/go-socket.io).

## Other Applications

I developed this software because i cannot make a "clean" implementation of
WebSockets using [WebDev](http://www.pcsoft.fr/webdev/index.html) because
it's based on Apache and IIS. So i needed an alternate server to handle the
websockets and allow other applications to talk easily to this server.

Communication is done via TCP using a simple framing mechanism.

## Inner Working

`webtalk` creates two servers, one for the websockets (called the Public
server), and another to communicate with another application (called the
Private server).

The Private server administration is done via a simple framing mechanism. Each
command consist of 5 text characters defining the message size, for example:

`00009[newuser]`

In this message, you see that there are 5 characters telling us that our
command will contain 9 characters. Then, our commands will have the format:

`[command] data`

In the example, the `[newuser]` doesn't have any data, it's just the command
with no parameters.

The response is exactly in the same format.

Why i didn't used another protocol? Because i wanted it to be easily
implemented in any tool. So far, i've implemented it on WinDev in less than
an hour, without using any external dependency, just plain W-Language.

The available commands are:

- `[newuser]` To create a new User, and returns: `[ok] ID!SALT`
- `[killuser] ID` Kill a User.
- `[whisper] ID!MESSAGE` Send a message to a specific user
- `[newchannel]` To create a new Channel, and returns: `[ok] ID`
- `[killchannel] ID` Kill a Channel.
- `[shout] ID!MESSAGE` Send a message to a entire channel.
- `[join] CHAN_ID!USER_ID` A user joins a Channel.
- `[leave] CHAN_ID!USER_ID` A user leaves a Channel.

Each application restart will clear all users and channels.

Also, your application must tell WebTalk to wait for a new user using the
`[newuser]` command. It will give you back a string in the next form `ID!SALT`,
which `ID` corresponds to your user Id. On the browser side, the user must send
the `ID!SALT` to the WebTalk server to login as that user. You cannot login two
browser to the same user id.

The responses from private server can be an `[error]` response, in wich case
it will contain the error message.

## License

(The MIT License)

Copyright (c) 2011 Hilario Pérez Corona &lt;hpcorona@gmail.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
