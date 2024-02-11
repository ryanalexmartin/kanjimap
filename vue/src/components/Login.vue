<template>
    <div>
        <h2>Login</h2>
        <input v-model="username" type="text" placeholder="Username">
        <input v-model="password" type="password" placeholder="Password">
        <button @click="login">Login</button>
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
            // Handle login here
        }
    }
}
</script>