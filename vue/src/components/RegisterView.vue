<!-- TODO write tests for this file -->
<template>
  <div class="center">
    <div class="register-view">
      <h1>Register</h1>
      <form @submit.prevent="register">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" required />

        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" required />

        <label for="email">Email (optional)</label>
        <input type="email" id="email" v-model="email" />

        <button type="submit">Register</button>
        <button type="button" @click="$emit('registered')">Cancel</button>
      </form>
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        username: '',
        password: '',
        email: ''
      }
    },
    methods: {
      register() {
        console.log('Registering user:', this.username);
        const apiUrl = process.env.VUE_APP_API_URL || 'http://localhost';
        const apiPort = process.env.VUE_APP_API_PORT || '8081';
        const url = `${apiUrl}:${apiPort}/register`;
        console.log('Sending request to:', url);

        fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: new URLSearchParams({
            username: this.username,
            password: this.password,
            email: this.email,
          }),
        })
          .then(response => {
            console.log('Response status:', response.status);
            console.log('Response headers:', response.headers);
            if (!response.ok) {
              return response.text().then(text => {
                throw new Error(`HTTP error! status: ${response.status}, message: ${text}`);
              });
            }
            return response.text();
          })
          .then(message => {
            console.log('Registration successful:', message);
            this.$emit('registered');
            this.login();
          })
          .catch(error => {
            console.error('There was a problem with the registration request:', error);
          });
      },
      login() {
        fetch(`${process.env.VUE_APP_API_URL}:${process.env.VUE_APP_API_PORT}/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: new URLSearchParams({
            username: this.username,
            password: this.password,
          }),
        })
          .then(response => {
            if (!response.ok) {
              throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text();
          })
          .then(message => {
            console.log(message);
            // Handle successful login here, e.g. by changing Vue state or redirecting
          })
          .catch(error => {
            console.log('There was a problem with the login request.', error);
          });
      }

    }
  }
</script>

<style scoped>
.register-view {
  align-self: center;
  margin: 20 auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  width: 300px;
}

input {
  display: block;
  margin: 10px;
  width: 90%;
  padding: 10px;
  font-size: 16px;
}

.center {
  display: flex;
  justify-content: center;
  align-items: center;
  /* height: 100vh; */
  padding-bottom: 20px;
}
</style>
