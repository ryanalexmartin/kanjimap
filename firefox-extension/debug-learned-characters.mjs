import fetch from 'node-fetch';

const API_URL = 'http://localhost:8081';
const username = 'insomagent';
const password = 'NihonG0!';

async function debugLearnedCharacters() {
    try {
        // Step 1: Login to get the token
        const loginResponse = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
        });

        if (!loginResponse.ok) {
            throw new Error(`Login failed with status ${loginResponse.status}: ${await loginResponse.text()}`);
        }

        const loginData = await loginResponse.json();
        const token = loginData.token;

        console.log('Login successful. Token:', token);

        // Step 2: Fetch learned characters
        const charactersUrl = `${API_URL}/learned-characters?username=${encodeURIComponent(username)}`;
        console.log('Fetching learned characters from:', charactersUrl);
        
        const charactersResponse = await fetch(charactersUrl, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!charactersResponse.ok) {
            const responseText = await charactersResponse.text();
            throw new Error(`Fetching learned characters failed with status ${charactersResponse.status}: ${responseText}`);
        }

        const learnedCharacters = await charactersResponse.json();
        console.log('Learned characters:', learnedCharacters);

    } catch (error) {
        console.error('Debug error:', error);
    }
}

debugLearnedCharacters();
