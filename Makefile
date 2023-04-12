# IMPORTANT: saving tg token is highly unsecure thing.
# DON'T do that in the real life!
# Even if it's created with learning motives like this one.

token := $(shell echo $(TG_TOKEN))

start:
	go build && ./read-adviser-bot -tg_bot_token '${token}'

start_sqlite:
	go build && ./read-adviser-bot -tg_bot_token '${token}' -storage_type 'sqlite'

start_files:
	go build && ./read-adviser-bot -tg_bot_token '${token}' -storage_type 'files'

.PHONY: start start_sqlite start_files