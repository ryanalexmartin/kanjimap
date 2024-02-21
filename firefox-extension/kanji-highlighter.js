async function loadKanjiDictionary() {
  const response = await fetch(browser.runtime.getURL('variant-WordData.json'));
  return response.json();
}

const kanjiDictionaryPromise = loadKanjiDictionary(); 

// save entire kanji dictionary to local storage
// localStorage.setItem('kanjiDictionary', JSON.stringify(kanjiDictionary));  // TODO Maybe if we did this it would be faster.

async function ensureToken() {
  let token = localStorage.getItem('token');
  if (!token) {
    console.log('No token found, fetching a new one...');
    const username = localStorage.getItem('username');
    const password = localStorage.getItem('password');
    try {
      const response = await fetch('https://kanjimap.cargocult.tech/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `username=${username}&password=${password}`
      });
      const data = await response.json();
      if (data.error) {
        console.error('Error fetching token:', data.error);
      } else {
        localStorage.setItem('token', data.token);
        console.log('Token saved:', data.token);
        token = data.token;
      }
    } catch (error) {
      console.error('Failed to fetch token:', error);
    }
  }
  return token; 
}

async function getZhuyin(kanji) {
  const kanjiDictionary = await kanjiDictionaryPromise; 
  const entry = kanjiDictionary.find(entry => entry.word === kanji);
  if (entry) {
    return entry.meanings[0].bopomofo;
  } else {
    return ''; 
  }
}

async function getLearnedKanji() {
  const token = await ensureToken(); 
  if (token) {
    const username = 'ryan'; // TODO: Get username dynamically if needed
    try {
      const response = await fetch(`https://kanjimap.cargocult.tech/fetch-characters?username=${username}`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await response.json();
      return data.map(kanji => kanji.character);
    } catch (error) {
      console.error('Failed to fetch learned kanji:', error);
      return []; 
    }
  } else {
    return []; 
  }
}

async function highlightKanjiInTextNode(node, learnedKanji) {
  const kanjiRegex = /[\u4e00-\u9faf\u3400-\u4dbf]/g;

  if (node.nodeType === Node.TEXT_NODE && kanjiRegex.test(node.nodeValue)) {
    const text = node.nodeValue;
    const fragment = document.createDocumentFragment();
    let lastIdx = 0;

    // Reset lastIndex because test method updates it
    kanjiRegex.lastIndex = 0;

    let match;
    while ((match = kanjiRegex.exec(text)) !== null) {
      // Append text before kanji

      fragment.appendChild(document.createTextNode(text.slice(lastIdx, match.index)));
        if (!learnedKanji.includes(match[0])) {
          const zhuyin = await getZhuyin(match[0]);
          const ruby = document.createElement('ruby');
          ruby.appendChild(document.createTextNode(match[0]));
          const rt = document.createElement('rt');
          rt.appendChild(document.createTextNode(zhuyin));
          ruby.appendChild(rt);
          fragment.appendChild(ruby);
        } else {
          // Append the kanji as is if it's already learned
          fragment.appendChild(document.createTextNode(match[0]));
        }


      lastIdx = match.index + match[0].length;
    }

    // Append any remaining text
    if (lastIdx < text.length) {
      fragment.appendChild(document.createTextNode(text.slice(lastIdx)));
    }
    node.parentNode.replaceChild(fragment, node);
  } else if (node.nodeType === Node.ELEMENT_NODE && !["SCRIPT", "STYLE", "RUBY"].includes(node.tagName.toUpperCase())) {
    // Clone NodeList to an array to avoid live collection issues during async operations
    Array.from(node.childNodes).forEach(child => highlightKanjiInTextNode(child, learnedKanji));
  }
}

async function highlightKanji() {
  try {
    const learnedKanji = await getLearnedKanji();
    highlightKanjiInTextNode(document.body, learnedKanji);
  } catch (error) {
    console.error('Failed to highlight kanji:', error);
  }
}

async function main() {
  await highlightKanji();
  const observer = new MutationObserver(highlightKanji);
  observer.observe(document.body, { childList: true, subtree: true });
}

// I want to be able to toggle the extension on and off from the icon, haven't figured out how to do that yet
// browser.runtime.onMessage.addListener(function(request, sender, sendResponse) {
//   if (request.extensionEnabled) {
//     main();
//   }
// });

main();

