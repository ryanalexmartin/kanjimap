let learnedCharacters = new Set();
let zhuyinData = new Map();

async function loadZhuyinData() {
  try {
    const response = await fetch(browser.runtime.getURL('characters.json'));
    const data = await response.json();
    data.forEach(entry => {
      if (entry.meanings && entry.meanings.length > 0) {
        zhuyinData.set(entry.word, entry.meanings[0].bopomofo);
      }
    });
    console.log('Zhuyin data loaded successfully');
  } catch (error) {
    console.error('Failed to load Zhuyin data:', error);
  }
}

function getZhuyin(kanji) {
  return zhuyinData.get(kanji) || '';
}

function injectStyles() {
  const style = document.createElement('style');
  style.textContent = `
    .kanji-highlight {
      display: inline-flex;
      flex-direction: row;
      align-items: center;
      vertical-align: middle;
      margin-right: 0.2em;
    }
    .kanji-highlight ruby {
      display: inline-flex;
      flex-direction: row;
      align-items: flex-start;
      text-align: center;
    }
    .kanji-highlight rb {
      display: inline-block;
      font-size: 1em;
      line-height: 1;
    }
    .kanji-highlight rt {
      display: inline-block;
      writing-mode: vertical-rl;
      text-orientation: upright;
      font-size: 0.4em;
      line-height: 1;
      text-align: start;
      margin-left: 0.1em;
      color: #666;
      font-weight: normal;
      max-height: 1.5em; /* Adjust based on your needs */
    }
  `;
  document.head.appendChild(style);
}

function highlightKanji(node) {
  if (node.nodeType === Node.TEXT_NODE) {
    const text = node.textContent;
    const kanjiRegex = /[\u4e00-\u9faf]/g;
    let match;
    let lastIndex = 0;
    const fragments = [];

    while ((match = kanjiRegex.exec(text)) !== null) {
      if (match.index > lastIndex) {
        fragments.push(document.createTextNode(text.slice(lastIndex, match.index)));
      }

      const kanji = match[0];
      if (!learnedCharacters.has(kanji)) {
        const span = document.createElement('span');
        span.className = 'kanji-highlight';
        const ruby = document.createElement('ruby');
        const rb = document.createElement('rb');
        rb.textContent = kanji;
        ruby.appendChild(rb);
        const rt = document.createElement('rt');
        rt.textContent = getZhuyin(kanji);
        ruby.appendChild(rt);
        span.appendChild(ruby);
        fragments.push(span);
      } else {
        fragments.push(document.createTextNode(kanji));
      }

      lastIndex = kanjiRegex.lastIndex;
    }

    if (lastIndex < text.length) {
      fragments.push(document.createTextNode(text.slice(lastIndex)));
    }

    if (fragments.length > 1) {
      const parent = node.parentNode;
      const container = document.createElement('span');
      fragments.forEach(fragment => container.appendChild(fragment));
      parent.replaceChild(container, node);
    }
  } else if (node.nodeType === Node.ELEMENT_NODE && !['SCRIPT', 'STYLE', 'TEXTAREA', 'RUBY', 'RT', 'RB'].includes(node.tagName)) {
    Array.from(node.childNodes).forEach(highlightKanji);
  }
}


function updateHighlights() {
  console.log('Updating highlights');
  browser.runtime.sendMessage({ action: 'getLearnedCharacters' })
    .then(response => {
      if (response && response.learnedCharacters) {
        learnedCharacters = new Set(response.learnedCharacters);
        requestAnimationFrame(() => highlightKanji(document.body));
      } else {
        console.error('Failed to get learned characters');
      }
    })
    .catch(error => {
      console.error('Error updating highlights:', error);
    });
}

// Load Zhuyin data and inject styles
loadZhuyinData().then(() => {
  injectStyles();
  updateHighlights();
}).catch(error => {
  console.error('Error during initialization:', error);
});

// Listen for changes in the DOM
const observer = new MutationObserver(mutations => {
  mutations.forEach(mutation => {
    mutation.addedNodes.forEach(node => {
      if (node.nodeType === Node.ELEMENT_NODE) {
        requestAnimationFrame(() => highlightKanji(node));
      }
    });
  });
});

observer.observe(document.body, { childList: true, subtree: true });

// Update highlights every 5 minutes
setInterval(updateHighlights, 5 * 60 * 1000);

console.log('Content script loaded');
