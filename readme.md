## DaaC2 - Discord as a C2

I wrote this because I was bored and wanted to mess around with a Discord bot, see if I could build a line of communication between a C2 server and an agent via Discord. Turns out it works pretty well. Ideally, with more time, I would use this as an initial dropper to then upgrade my shell via injection. But you know, I'm not a red teamer - so who knows what I'd actually do with it. I just like making things.

More information can be seen here: https://crawl3r.github.io/2020-01-25/DaaC2

Usage:
* Create an app & bot in discord dev portal, add the bot to a new server/channel created by you
* Obtain auth token for bot in dev portal, set the Token value in c2discord.go
* In addition to the Token, a channel ID is required - I obtained this by debugging the first message the agent received. I honestly don't know if there's another way, but when you have that, add it to the c2discord.go file with the token and you're good to go
* Build and run the C2 server. Should support both Mac OSX and Linux - however it's only been properly tested on Mac OSX
* Build and deploy agent (Should work on Mac OSX, Windows and Linux)
* Start using the server Cli to control agent (type 'help' for a quick drop down of initial commands)

Don't use this for anything illegal. I haven't invested any time into making sure it can hide properly or be at all stealthy, so chances are you'll get pwned before you pwn them :) I am not responsible for anything you do with this, so just use your brain - please. A skill not many humans have anymore.

Development will be on-going but slow, so feel free to tweet me any comments etc @monobehaviour

Known issues: 
* shellcode injection doesn't work on Windows (w.i.p)
* C2 can't be built for Windows (does anyone actually want that?)