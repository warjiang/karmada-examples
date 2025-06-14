from dotenv import load_dotenv
from telegram.ext import ApplicationBuilder
import asyncio
import os
from datetime import datetime

async def main():
    load_dotenv()
    bot_token = os.environ.get('BOT_TOKEN')
    chat_id = os.environ.get('CHAT_ID')
    application = ApplicationBuilder().token(bot_token).build()
    current_date = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    msg = f'{current_date} image sync finished'
    await application.bot.send_message(chat_id, msg)


if __name__ == '__main__':
    asyncio.run(main())