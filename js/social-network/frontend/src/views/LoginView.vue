<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Welcome back</h1>
      <p class="auth-sub">Sign in to your account</p>

      <div v-if="sessionMessage" class="error-msg" style="margin-bottom:16px;background:var(--accent-bg);border-color:var(--border2);color:var(--accent)">
        {{ sessionMessage }}
      </div>

      <div class="form-field">
        <label>Email</label>
        <input class="input" v-model="email" type="email" placeholder="you@example.com" @keyup.enter="login" />
      </div>
      <div class="form-field">
        <label>Password</label>
        <input class="input" v-model="password" type="password" placeholder="••••••••" @keyup.enter="login" />
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>

      <button class="btn btn-primary w-full" style="margin-top:20px" :disabled="loading" @click="login">
        {{ loading ? 'Signing in…' : 'Sign in' }}
      </button>
      <p class="auth-footer" style="margin-top:20px">
        No account? <router-link to="/register">Create one</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '../services/api'
import { setLoggedIn } from '../stores/auth'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const sessionMessage = ref('')
const router = useRouter()

onMounted(() => {
  const msg = sessionStorage.getItem('session_msg')
  if (msg) { sessionMessage.value = msg; sessionStorage.removeItem('session_msg') }
})
async function login() {
  if (!email.value || !password.value) { error.value = 'Email and password are required.'; return }
  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRe.test(email.value.trim())) { error.value = 'Please enter a valid email address.'; return }
  loading.value = true; error.value = ''
  try {
    await apiFetch('/login', { method: 'POST', body: JSON.stringify({ email: email.value, password: password.value }) })
    setLoggedIn(true)
    router.push('/feed')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
