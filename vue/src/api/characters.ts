import axios from 'axios'
import { Character } from '@/types/types'

const API_URL = import.meta.env.VITE_API_BASE_URL

export async function fetchCharacters(username: string): Promise<Character[]> {
  try {
    console.log('API_URL', API_URL)
    const response = await axios.get(
      `${API_URL}/fetch-characters`,
      {
        params: { username },
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      }
    );
    
    // Check if the response is HTML
    if (typeof response.data === 'string' && response.data.trim().startsWith('<!DOCTYPE html>')) {
      console.error('Received HTML response instead of character data');
      throw new Error('Invalid response format: received HTML');
    }
    
    console.log('API response:', response.data);
    
    // Validate that the response is an array
    if (!Array.isArray(response.data)) {
      console.error('Invalid response format: expected an array');
      throw new Error('Invalid response format: expected an array');
    }
    
    return response.data;
  } catch (error) {
    console.error('Error fetching characters:', error);
    throw error;
  }
}

export async function updateCharacterLearned(character: Character): Promise<Character> {
  const response = await axios.post(`${API_URL}/learn-character`, {
    username: localStorage.getItem('username'),
    chinese_character: character.char,
    pinyin: character.pinyin,
    characterId: character.id,
    learned: !character.learned
  }, {
    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
  })
  return {
    ...character,
    learned: response.data.learned
  }
}