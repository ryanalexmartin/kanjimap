<template>
    <div>
        <h2>Login</h2>
        <input v-model="username" type="text" placeholder="Username">
        <input v-model="password" type="password" placeholder="Password">
        <button @click="login">Login</button>
        <button @click="$emit('register')">Register</button>
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
            fetch('http://localhost:8081/login', {
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
h2 {
    margin-bottom: 20px;
}
input {
    display: block;
    margin-bottom: 10px;
    width: 100%;
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
</style>