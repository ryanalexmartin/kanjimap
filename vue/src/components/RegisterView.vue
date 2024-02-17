<!-- TODO write tests for this file -->
<template>
    <div class="center">
        <div class="register-view">
            <h1>Register</h1>
            <form @submit.prevent="register">
                <label for="username">Username</label>
                <input type="text" id="username" v-model="username" required>

                <label for="password">Password</label>
                <input type="password" id="password" v-model="password" required>

                <label for="email">Email (optional)</label>
                <input type="email" id="email" v-model="email">

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
            fetch('http://localhost:8081/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams({
                    username: this.username,
                    password: this.password,
                    email: this.email
                }),
            })
                .then(response => {
                    if (!response.ok) {
                        alert(response.statusText);
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
