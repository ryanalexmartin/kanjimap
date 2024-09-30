let authToken = '';
let learnedCharacters = [];

async function authenticate() {
  const storage = await browser.storage.local.get(['username', 'password']);
  const { username, password } = storage;
  
  if (!username || !password) {
    console.error('Username or password not set');
    return false;
  }

  try {
    const response = await fetch('http://localhost:8081/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `username=${username}&password=${password}`
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    authToken = data.token;
    return true;
  } catch (error) {
    console.error('Authentication failed:', error);
    return false;
  }
}

async function fetchLearnedCharacters() {
  if (!authToken) {
    const authenticated = await authenticate();
    if (!authenticated) return;
  }

  try {
    const { username } = await browser.storage.local.get('username');
    const response = await fetch(`http://localhost:8081/learned-characters?username=${username}`, {
      headers: { 'Authorization': `Bearer ${authToken}` }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    learnedCharacters = await response.json();
    await browser.storage.local.set({ learnedCharacters });
  } catch (error) {
    console.error('Failed to fetch learned characters:', error);
  }
}

browser.runtime.onMessage.addListener((request, sender, sendResponse) => {
  console.log('Received message in background:', request);
  if (request.action === 'getLearnedCharacters') {
    fetchLearnedCharacters().then(() => {
      console.log('Sending learned characters:', learnedCharacters);
      sendResponse({ learnedCharacters: learnedCharacters });
    });
    return true; // Indicates that we will send a response asynchronously
  } else if (request.action === 'updateSettings') {
    // Relay the settings update to all content scripts
    browser.tabs.query({}).then((tabs) => {
      for (let tab of tabs) {
        browser.tabs.sendMessage(tab.id, {
          action: 'updateSettings',
          settings: request.settings
        }).catch((error) => console.log(`Error sending message to tab ${tab.id}:`, error));
      }
    });
    sendResponse({ status: 'Settings update relayed to content scripts' });
    return true;
  }
});


// Initial fetch on extension load
fetchLearnedCharacters();

console.log('Background script loaded');
