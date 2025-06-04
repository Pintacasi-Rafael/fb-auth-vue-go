<template>
  <div class="landing-page">
    <h1>Welcome, {{ userName }}!</h1>
    <p>Your email: {{ userEmail }}</p>
    <p>Your JWT token:</p>
    <textarea readonly rows="5" cols="50" v-model="jwtToken"></textarea>
    <button @click="logout">Logout</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const jwtToken = ref('')
const userName = ref('')
const userEmail = ref('')

const router = useRouter()

function parseJwt(token) {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join(''),
    )
    return JSON.parse(jsonPayload)
  } catch {
    return null
  }
}

onMounted(() => {
  const token = localStorage.getItem('jwtToken') || ''
  if (!token) {
    router.push('/')
    return
  }
  jwtToken.value = token

  const payload = parseJwt(token)
  if (payload) {
    userName.value = payload.name || ''
    userEmail.value = payload.email || ''
  }
})

function logout() {
  localStorage.removeItem('jwtToken')
  router.push('/')
}
</script>

<style scoped>
.landing-page {
  max-width: 600px;
  margin: 30px auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: Arial, sans-serif;
  text-align: center;
}
textarea {
  width: 100%;
  margin: 10px 0;
  font-family: monospace;
}
button {
  padding: 10px 20px;
  background-color: #4267b2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
button:hover {
  background-color: #365899;
}
</style>
