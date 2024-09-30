document.addEventListener('DOMContentLoaded', function() {
  const loginForm = document.getElementById('loginForm');
  const loginBtn = document.getElementById('loginBtn');
  const logoutBtn = document.getElementById('logoutBtn');
  const status = document.getElementById('status');
  const onOffSwitch = document.getElementById('onOffSwitch');
  const onOffSwitchDiv = document.getElementById('onOffSwitchDiv');
  const orientationSelect = document.getElementById('orientationSelect');
  const orientationSelectDiv = document.getElementById('orientationSelectDiv');

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

  orientationSelect.addEventListener('change', (event) => {
    const orientation = event.target.value;
    updateOrientation(orientation);
  });

  function updateOrientation(orientation) {
    browser.storage.local.set({ rubyOrientation: orientation });
    browser.tabs.query({active: true, currentWindow: true}).then((tabs) => {
      browser.tabs.sendMessage(tabs[0].id, {
        action: 'updateSettings',
        settings: { rubyOrientation: orientation }
      });
    });
  }

  // Load saved orientation
  browser.storage.local.get('rubyOrientation').then(result => {
    console.log(result)
    if (result.rubyOrientation) {
      document.getElementById('orientationSelect').value = result.rubyOrientation;
    }
  });


  // Load saved zhuyinEnabled and set checkbox state
  browser.storage.local.get('zhuyinEnabled').then(result => {
    const zhuyinEnabled = result.zhuyinEnabled === 'on';
    onOffSwitch.checked = zhuyinEnabled;
    updateZhuyinState(zhuyinEnabled);
  });

  onOffSwitch.addEventListener('change', (event) => {
    const zhuyinEnabled = event.target.checked;
    updateZhuyinState(zhuyinEnabled);
  });

  function updateZhuyinState(enabled) {
    const zhuyinState = enabled ? 'on' : 'off';
    browser.storage.local.set({ zhuyinEnabled: zhuyinState });
    browser.tabs.query({active: true, currentWindow: true}).then((tabs) => {
      browser.tabs.sendMessage(tabs[0].id, {
        action: 'updateSettings',
        settings: { zhuyinEnabled: zhuyinState }
      });
    });
  }



  logoutBtn.addEventListener('click', function() {
    browser.storage.local.remove(['authToken', 'username', 'password']).then(() => {
      showLoggedOutUI();
      status.textContent = 'Logged out successfully!';
    });
  });

  function showLoggedInUI() {
    loginForm.style.display = 'none';
    logoutBtn.style.display = 'block';
    onOffSwitch.style.display = 'block';
    onOffSwitchDiv.style.display = 'block';
    orientationSelect.style.display = 'block';
    orientationSelectDiv.style.display = 'block';
  }

  function showLoggedOutUI() {
    loginForm.style.display = 'block';
    logoutBtn.style.display = 'none';
    onOffSwitch.style.display = 'none';
    onOffSwitchDiv.style.display = 'none';
    orientationSelect.style.display = 'none';
    orientationSelectDiv.style.display = 'none';
  }
});
