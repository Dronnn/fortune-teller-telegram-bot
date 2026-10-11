import { positiveAnswers, uncertainAnswers, negativeAnswers, mysticalEmoji } from './answers.js';

function pick(list) {
  return list[Math.floor(Math.random() * list.length)];
}

export function shake() {
  const roll = Math.floor(Math.random() * 100);
  const pool = roll < 50 ? positiveAnswers : roll < 75 ? uncertainAnswers : negativeAnswers;
  return `${pick(mysticalEmoji)} ${pick(pool)}`;
}
