# IMPORTANT: saving tg token is highly unsecure thing.
# DON'T do that in the real life!
# Even if it's created with learning motives like this one.

token := $(shell echo $(TG_TOKEN))

.PHONY start:
	go build && ./read-adviser-bot -tg-bot-token '${token}'

