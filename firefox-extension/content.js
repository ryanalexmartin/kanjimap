const RUBY_ORIENTATIONS = {
  VERTICAL_RIGHT: 'vertical-right',
  HORIZONTAL_ABOVE: 'horizontal-above',
  HORIZONTAL_BELOW: 'horizontal-below'
};
let learnedCharacters = new Set();
let zhuyinData = new Map();
let currentOrientation = RUBY_ORIENTATIONS.HORIZONTAL_ABOVE;
let zhuyinEnabled = 'off';

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
      display: inline-block;
      vertical-align: baseline;
      line-height: normal;
      position: relative;
    }
    .kanji-highlight ruby {
      display: inline-flex;
      vertical-align: baseline;
      line-height: 1;
    }
    .kanji-highlight rb {
      display: inline-block;
      font-size: 1em;
      line-height: inherit;
    }
    .kanji-highlight rt {
      display: inline-block;
      font-size: 0.3em;
      line-height: normal;
      text-align: start;
      color: #666;
      font-weight: normal;
    }
    .kanji-highlight.vertical-right {
      margin-right: 0.4em;
    }
    .kanji-highlight.vertical-right rt {
      display:flex;
      justify-content:center;
      align-items:center;
      writing-mode: vertical-rl;
      text-orientation: upright;
      position: absolute;
      top: 0;
      right: -1.2em;
      height: 100%;
    }
    .kanji-highlight.horizontal-above rt {
      position: absolute;
      top: -0.6em;
      left: 50%;
      transform: translateX(-50%);
      white-space: nowrap;
    }
    .kanji-highlight.horizontal-below rt {
      position: absolute;
      bottom: -0.9em;
      left: 50%;
      transform: translateX(-50%);
      white-space: nowrap;
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
        span.className = `kanji-highlight ${currentOrientation}`;
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

function setRubyOrientation(orientation) {
  currentOrientation = orientation;
  browser.storage.local.set({ rubyOrientation: orientation });
  updateHighlights();
}

function setEnabledZhuyinHighlighting(value) {
  zhuyinEnabled = value;
  browser.storage.local.set({ zhuyinEnabled: value });
  updateHighlights();
}

function updateHighlights() {
  console.log('Updating highlights. Zhuyin enabled:', zhuyinEnabled);
  removeExistingHighlights();
  
  if (zhuyinEnabled === 'on') {
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
}


function removeExistingHighlights() {
  console.log('Removing existing highlights');
  const highlights = document.querySelectorAll('.kanji-highlight');
  highlights.forEach(highlight => {
    const parent = highlight.parentNode;
    while (highlight.firstChild) {
      parent.insertBefore(highlight.firstChild, highlight);
    }
    parent.removeChild(highlight);
  });
}

function applySettings(settings) {
  console.log('Applying settings:', settings);
  let settingsChanged = false;

  if (settings.rubyOrientation && settings.rubyOrientation !== currentOrientation) {
    currentOrientation = settings.rubyOrientation;
    settingsChanged = true;
    console.log('Updated orientation to:', currentOrientation);
  }

  if (settings.zhuyinEnabled !== undefined && settings.zhuyinEnabled !== zhuyinEnabled) {
    zhuyinEnabled = settings.zhuyinEnabled;
    settingsChanged = true;
    console.log('Updated Zhuyin enabled to:', zhuyinEnabled);
  }

  if (settingsChanged) {
    updateHighlights();
  }
}

// Listen for messages from popup
browser.runtime.onMessage.addListener((message, sender, sendResponse) => {
  console.log('Received message:', message);
  if (message.action === 'updateSettings') {
    applySettings(message.settings);
    sendResponse({status: 'Settings updated'});
  }
  return true;
});

// Load saved settings and initialize
browser.storage.local.get(['rubyOrientation', 'zhuyinEnabled']).then(result => {
  console.log('Loaded settings:', result);
  if (result.rubyOrientation) {
    currentOrientation = result.rubyOrientation;
  }
  if (result.zhuyinEnabled) {
    zhuyinEnabled = result.zhuyinEnabled;
  }
  loadZhuyinData().then(() => {
    injectStyles();
    updateHighlights();
  }).catch(error => {
    console.error('Error during initialization:', error);
  });
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
