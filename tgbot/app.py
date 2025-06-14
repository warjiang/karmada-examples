from dotenv import load_dotenv
from telegram.ext import ApplicationBuilder
import asyncio
import os
from datetime import datetime
import sys

async def main():
    args = sys.argv
    for arg in args:
        print(arg)
    if len(args) < 2:
        print('Usage: uv run app.py <msg>')
        exit(1)

    msg = sys.argv[1]

    load_dotenv()
    bot_token = os.environ.get('BOT_TOKEN')
    chat_id = os.environ.get('CHAT_ID')
    application = ApplicationBuilder().token(bot_token).build()
    current_date = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    msg_with_time = f'{current_date} {msg}'
    await application.bot.send_message(chat_id, msg_with_time)


if __name__ == '__main__':
    asyncio.run(main())
