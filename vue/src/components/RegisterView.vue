<template>
    <div class="register-view">
        <h1>Register</h1>
        <form @submit.prevent="register">
            <label for="username">Username:</label>
            <input type="text" id="username" v-model="username" required>

            <label for="password">Password:</label>
            <input type="password" id="password" v-model="password" required>

            <button type="submit">Register</button>
        </form>
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
        register() {
            console.log('Registering user:', this.username);
            fetch('http://localhost:8081/register', {
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
                    // Emit an event to the parent component
                    this.$emit('registered');
                    this.login();
                })
                .catch(error => {
                    console.log('There was a problem with the registration request.', error);
                });
        },
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
    max-width: 300px;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
</style>