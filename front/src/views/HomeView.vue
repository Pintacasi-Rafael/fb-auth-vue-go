<template>
  <div>
    <button @click="redirectToFacebook">Login with Facebook</button>
    <p v-if="error" style="color: red">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const error = ref('')
const jwtToken = ref('')

const fbAppID = '1245020367235225'
const redirectUri = 'http://localhost:8080/auth/facebook/callback' // backend redirect URI

const router = useRouter()
const route = useRoute()

function redirectToFacebook() {
  const fbOAuthUrl = new URL('https://www.facebook.com/v17.0/dialog/oauth')
  fbOAuthUrl.searchParams.set('client_id', fbAppID)
  fbOAuthUrl.searchParams.set('redirect_uri', redirectUri)
  fbOAuthUrl.searchParams.set('scope', 'email')
  fbOAuthUrl.searchParams.set('response_type', 'code')

  // Redirect user to Facebook OAuth page
  window.location.href = fbOAuthUrl.toString()
}

async function exchangeCodeForToken(code) {
  error.value = ''
  try {
    const res = await fetch(
      'http://localhost:8080/auth/facebook/callback?code=' + encodeURIComponent(code),
    )
    if (!res.ok) {
      const errData = await res.json()
      error.value = errData.error || 'Failed to authenticate'
      return
    }
    const data = await res.json()
    jwtToken.value = data.token
    localStorage.setItem('jwtToken', data.token)
    router.push('/landing')
  } catch (err) {
    error.value = 'Network error: ' + err.message
  }
}

onMounted(() => {
  // Handle token passed via URL query param after redirect from backend
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')
  if (token) {
    localStorage.setItem('jwtToken', token)
    router.push('/landing')

    // Clean URL to remove token from address bar
    window.history.replaceState({}, document.title, window.location.pathname)
    return
  }

  // Handle session expiration message
  if (route.query.sessionExpired) {
    error.value = 'Your session has expired. Please log in again.'
  }

  // Legacy code for code-based OAuth flow (optional, if still used)
  const code = urlParams.get('code')
  if (code) {
    window.history.replaceState({}, document.title, window.location.pathname)
    exchangeCodeForToken(code)
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
