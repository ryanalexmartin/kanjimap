<!-- TODO write tests for this file -->
<template>
    <div class="center">
        <div class="login-view">
            <h2>Login</h2>
            <input v-model="username" type="text" placeholder="Username">
            <input v-model="password" type="password" placeholder="Password">
            <button @click="login">Login</button>
            <button @click="$emit('register')">Register</button>
        </div>
    </div>
</template>

<script>

export default {
    data() {
        return {
            username: '',
            password: ''
        }
    },
    methods: {
        login() {
            console.log('Logging in user:', this.username);
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
                    return response.json();
                })
                .then(response => {
                    // save the token in local storage
                    const token = response.token;
                    console.log('token', token);
                    localStorage.setItem('token', token);
                    localStorage.setItem('username', this.username);
                    this.$emit('login', this.username);
                })
                .catch(error => {
                    console.log('There was a problem with the login request.', error);
                });
        }
    }
}
</script>

<style scoped>
.login-view {
    align-self: center;
    margin: 20 auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    width: 300px;
}

h2 {
    margin-bottom: 20px;
}
input {
    display: block;
    margin-bottom: 10px;
    width: 90%;
    padding: 10px;
    font-size: 16px;
}
button {
    padding: 10px;
    margin: 10px;
    font-size: 16px;
    background-color: #4CAF50;
    color: white;
    border: none;
    cursor: pointer;
}
button:hover {
    background-color: #45a049;
}
.center {
    display: flex;
    justify-content: center;
    align-items: center;
    /* height: 100vh; */
    padding-bottom: 20px;
}
</style>
