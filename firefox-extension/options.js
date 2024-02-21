document.getElementById('options-form').addEventListener('submit', function(e) {
  e.preventDefault();
  var username = document.getElementById('username').value;
  var password = document.getElementById('password').value;
  // var enableZhuyin = document.getElementById('enable-zhuyin').checked;
  // var enableHighlighting = document.getElementById('enable-highlighting').checked;

    localStorage.setItem('username', username);
    localStorage.setItem('password', password);
    // localStorage.setItem('enableZhuyin', enableZhuyin);
    // localStorage.setItem('enableHighlighting', enableHighlighting);
    console.log('username and password saved: ' + username + ' ' + password);

    // Make a request to the kanjimap API to verify the username and password
    // URL encoded form data
    fetch('https://kanjimap.cargocult.tech/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: 'username=' + username + '&password=' + password
    })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      if (data.error) {
        document.getElementById('error').textContent = data.error;
      } else {
        localStorage.setItem('token', data.token);
        console.log('token saved: ' + data.token);
      }
    })
})

