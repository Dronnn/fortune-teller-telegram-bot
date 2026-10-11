import { api } from 'sdk';
import { shake } from '../lib/fortuneball.js';
import { hasTrigger } from '../lib/triggers.js';

// The go bot fetched its username once at startup; here we resolve it lazily
// and remember it for the life of the isolate (may reset between invocations).
let botMe = null;

async function getMe() {
  if (!botMe) botMe = await api.getMe();
  return botMe;
}

// go equivalent: msg.IsCommand() + msg.Command(), extended for group chats:
// a command typed after a leading @mention ("@bot /help") counts too — that is
// how people ask a bot for help in a group.
function getCommand(message) {
  const text = message.text;
  if (!text) return null;
  const entity = (message.entities ?? []).find((e) => e.type === 'bot_command');
  if (!entity) return null;
  const before = text.slice(0, entity.offset);
  if (entity.offset !== 0 && !/^@\S+\s*$/.test(before)) return null;
  return text.slice(entity.offset, entity.offset + entity.length).slice(1).split('@')[0].toLowerCase();
}

// go equivalent: strings.Contains(text, "@"+username), but exact via the full
// handle — in groups this only matches a real @botname mention.
async function isMentioned(message, loweredText) {
  if (!loweredText.includes('@')) return false;
  const { username } = await getMe();
  if (!username) return false;
  return loweredText.includes(`@${username.toLowerCase()}`);
}

// go equivalent: msg.ReplyToMessage.From.UserName == botUsername
// (compared by id here, which is stronger)
async function repliedToBot(message) {
  const reply = message.reply_to_message;
  if (!reply || !reply.from) return false;
  const { id } = await getMe();
  return reply.from.id === id;
}

async function shouldRespond(message) {
  const text = message.text;
  if (!text) return false;
  const lowered = text.toLowerCase();

  // Private chat: trigger phrase or mention
  if (message.chat.type === 'private') {
    return hasTrigger(lowered) || (await isMentioned(message, lowered));
  }

  // Group chat: @mention or reply to bot
  return (await isMentioned(message, lowered)) || (await repliedToBot(message));
}

async function handleStart(message) {
  await api.sendMessage({
    chat_id: message.chat.id,
    text: '🔮 Я — магический шар предсказаний!\n\nЗадай мне вопрос, и я отвечу.\nНапиши /help чтобы узнать подробнее.',
  });
}

async function handleHelp(message) {
  const { username } = await getMe();
  await api.sendMessage({
    chat_id: message.chat.id,
    parse_mode: 'Markdown',
    text: `🔮 *Как спросить предсказание:*

Начни вопрос с одной из фраз:
• "Что будет, если..."
• "Будет ли..."
• "Стоит ли мне..."
• "Предскажи..."
• "Суждено ли..."
• "Повезёт ли..."
• "О великий шар..."

В групповом чате — упомяни @${username} или ответь на моё сообщение.`,
  });
}

export default async function (message) {
  const command = getCommand(message);
  if (command === 'start') return handleStart(message);
  if (command === 'help') return handleHelp(message);

  if (!(await shouldRespond(message))) return;

  await api.sendMessage({
    chat_id: message.chat.id,
    text: shake(),
    reply_parameters: { message_id: message.message_id },
  });
}
