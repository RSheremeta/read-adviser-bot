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
3. In order to run the Telegram Bot service, just do the ```make start``` command in terminal

Hint: the `make` builds the binary and runs the service whenever you invoke it, but you can also use a pre-built one [here](https://github.com/RSheremeta/read-adviser-bot/tags)