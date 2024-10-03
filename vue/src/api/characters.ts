import axios from 'axios'
import { Character } from '@/types/types'

const API_URL = import.meta.env.VITE_API_URL

export async function fetchCharacters(username: string): Promise<Character[]> {
  const response = await axios.get(
    `${API_URL}/fetch-characters`,
    {
      params: { username },
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
    }
  );
  return response.data;
}

export async function updateCharacterLearned(character: Character): Promise<Character> {
  const response = await axios.post(`${API_URL}/learn-character`, {
    username: localStorage.getItem('username'),
    chinese_character: character.char,
    characterId: character.id,
    learned: !character.learned
  }, {
    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
  })
  return response.data
}