# Whatsapp -> Discord bot

This bot just reposts whatsapp messages from a group chat (based on chat name) into a discord channel.

To set it up, youd need a `.env` file in your working directory, which has the following items:

```sh
# Discord bot token
TOKEN=""
# Discord Channel ID
GUILD_CHANNEL=""
# The role to @mention when sending a message
GUILD_ROLE=""
# The group chat name
CHAT_NAME=""
```

The bot currently supports the following message types:
- Normal text messages
- Audio messages
- Document messages
- Image messages
- Video messages
- Location messages (including live location, but doesn't track it *live*, send the initial location only)

To make this work, youll have to setup golang, then run the following:

```sh
go install github.com/ShadiestGoat/Whatsapp-Discord-Bot@latest
```

After installing, create a `.env` file as described above, and you will be able to run this correctly. Please do note, that this software will create a `deviceStore.db` file in the current directory & it is highly advised to not leak that. 

If it can't auto login, you will be prompted with a qr code through your terminal. Scan it using whatsapp mobile
