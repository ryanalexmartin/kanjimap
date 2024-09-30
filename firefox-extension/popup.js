document.addEventListener('DOMContentLoaded', function() {
  const loginForm = document.getElementById('loginForm');
  const loginBtn = document.getElementById('loginBtn');
  const logoutBtn = document.getElementById('logoutBtn');
  const highlightBtn = document.getElementById('highlightBtn');
  const status = document.getElementById('status');

  // Check if user is already logged in
  browser.storage.local.get('authToken').then((result) => {
    if (result.authToken) {
      showLoggedInUI();
    }
  });

  loginBtn.addEventListener('click', function() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('http://localhost:8081/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
    })
      .then(response => response.json())
      .then(data => {
        if (data.token) {
          browser.storage.local.set({ authToken: data.token, username: username, password: password });
          showLoggedInUI();
          status.textContent = 'Logged in successfully!';
          browser.runtime.sendMessage({ action: 'fetchLearnedCharacters' });
        } else {
          status.textContent = 'Login failed. Please try again.';
        }
      })
      .catch(error => {
        status.textContent = 'An error occurred. Please try again.';
        console.error('Login error:', error);
      });
  });

  document.getElementById('orientation').addEventListener('change', (event) => {
    const orientation = event.target.value;
    browser.tabs.query({active: true, currentWindow: true}).then((tabs) => {
      browser.tabs.sendMessage(tabs[0].id, {
        action: 'setRubyOrientation',
        orientation: orientation
      });
    });
    browser.storage.local.set({ rubyOrientation: orientation });
  });

  // Load saved orientation
  browser.storage.local.get('rubyOrientation').then(result => {
    if (result.rubyOrientation) {
      document.getElementById('orientation').value = result.rubyOrientation;
    }
  });

  logoutBtn.addEventListener('click', function() {
    browser.storage.local.remove(['authToken', 'username', 'password']).then(() => {
      showLoggedOutUI();
      status.textContent = 'Logged out successfully!';
    });
  });

  highlightBtn.addEventListener('click', function() {
    browser.tabs.query({ active: true, currentWindow: true }).then((tabs) => {
      browser.tabs.sendMessage(tabs[0].id, { action: 'updateHighlights' });
      status.textContent = 'Highlighting updated!';
    });
  });

  function showLoggedInUI() {
    loginForm.style.display = 'none';
    logoutBtn.style.display = 'block';
    highlightBtn.style.display = 'block';
  }

  function showLoggedOutUI() {
    loginForm.style.display = 'block';
    logoutBtn.style.display = 'none';
    highlightBtn.style.display = 'none';
  }
});
