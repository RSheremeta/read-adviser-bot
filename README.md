# read-adviser-bot
A Telegram Bot which consumes/stores/gives the pages for you to read


---
### FAQ
Available commands:
- /start
- /help
- /rnd (pick a random link from the saved list)

**Important:**
Picking up links is a one-time action!
If you make the Bot to pick any link for you, this link will be removed from the storage immediately.

### Usage
1. Get (or retrieve created) the Telegram Bot token value.
2. Set it as an env variable called `TG_TOKEN`.
3. There are 3 ways in order to run the Telegram Bot service:
- ```make start-sqlite``` command in terminal (uses sqlite as a storage source) 
- ```make start-files``` command in terminal (uses Gob encoded in files as a storage source) 
- just ```make start```  (uses **sqlite** by default)

Hint: the `make` commands do build the binary and run the service whenever you invoke any of it, but you can also use a pre-built one [here](https://github.com/RSheremeta/read-adviser-bot/tags)

### Demo

![](https://github.com/RSheremeta/read-adviser-bot/blob/master/demo.gif)