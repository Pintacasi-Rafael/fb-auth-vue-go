<template>
  <div>
    <button @click="redirectToFacebook">Login with Facebook</button>
    <p v-if="error" style="color: red">{{ error }}</p>
    <p v-if="jwtToken">JWT Token: {{ jwtToken }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const error = ref('')
const jwtToken = ref('')

const fbAppID = '1245020367235225'
const redirectUri = 'http://localhost:8080/auth/facebook/callback' // backend redirect URI

const router = useRouter()

function redirectToFacebook() {
  const fbOAuthUrl = new URL('https://www.facebook.com/v17.0/dialog/oauth')
  fbOAuthUrl.searchParams.set('client_id', fbAppID)
  fbOAuthUrl.searchParams.set('redirect_uri', redirectUri)
  fbOAuthUrl.searchParams.set('scope', 'email')
  fbOAuthUrl.searchParams.set('response_type', 'code')

  // Redirect user to Facebook OAuth page
  window.location.href = fbOAuthUrl.toString()
}

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')

  if (token) {
    jwtToken.value = token
    localStorage.setItem('jwtToken', token)
    window.history.replaceState({}, document.title, window.location.pathname)
    router.push('/landing')
  }
})
</script>

<style scoped>
button {
  padding: 10px 20px;
  background-color: #4267b2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}
button:hover {
  background-color: #365899;
}
</style>
