let extensionEnabled = true;

function toggleExtension(tab) {
  // Toggle the state
  extensionEnabled = !extensionEnabled;

  // Send a message to the content script with the new state
  browser.tabs.sendMessage(tab.id, {extensionEnabled: extensionEnabled});
}

// Listen for clicks on the extension icon
browser.browserAction.onClicked.addListener(toggleExtension);

