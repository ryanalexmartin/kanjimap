let learnedCharacters = [];

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
      if (!learnedCharacters.includes(kanji)) {
        const span = document.createElement('span');
        span.textContent = kanji;
        span.style.backgroundColor = 'yellow';
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
      fragments.forEach(fragment => parent.insertBefore(fragment, node));
      parent.removeChild(node);
    }
  } else if (node.nodeType === Node.ELEMENT_NODE && !['SCRIPT', 'STYLE', 'TEXTAREA'].includes(node.tagName)) {
    Array.from(node.childNodes).forEach(highlightKanji);
  }
}

function updateHighlights() {
  console.log('Updating highlights');
  browser.runtime.sendMessage({ action: 'getLearnedCharacters' })
    .then(response => {
      console.log('Received response:', response);
      if (response && response.learnedCharacters) {
        learnedCharacters = response.learnedCharacters;
        highlightKanji(document.body);
      } else {
        console.error('Failed to get learned characters');
      }
    })
    .catch(error => {
      console.error('Error updating highlights:', error);
    });
}

// Initial highlighting
updateHighlights();

// Listen for changes in the DOM
const observer = new MutationObserver(mutations => {
  mutations.forEach(mutation => {
    mutation.addedNodes.forEach(highlightKanji);
  });
});

observer.observe(document.body, { childList: true, subtree: true });

// Update highlights every 5 minutes
setInterval(updateHighlights, 5 * 60 * 1000);

console.log('Content script loaded');
