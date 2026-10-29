<template>
  <div class="auth-page">
    <div class="auth-card" style="max-width:500px">
      <h1>Create account</h1>
      <p class="auth-sub">Join the community today</p>

      <!-- Avatar picker -->
      <div style="display:flex;align-items:center;gap:16px;margin-bottom:20px;padding:16px;background:var(--surface2);border-radius:12px;border:1px solid var(--border)">
        <div style="width:64px;height:64px;border-radius:50%;overflow:hidden;flex-shrink:0;background:var(--surface3);display:flex;align-items:center;justify-content:center">
          <img v-if="avatarPreview" :src="avatarPreview" style="width:100%;height:100%;object-fit:cover" />
          <svg v-else style="width:28px;height:28px;color:var(--text-3)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div>
          <div style="font-size:0.85rem;font-weight:600;color:var(--text);margin-bottom:6px">Profile Photo <span style="color:var(--text-3);font-weight:400">(optional)</span></div>
          <label style="cursor:pointer">
            <input type="file" accept=".jpg,.jpeg,.png,.gif" style="display:none" ref="avatarInput" @change="pickAvatar" />
            <span class="btn btn-ghost btn-sm" style="display:inline-flex;align-items:center">
              <svg style="width:14px;height:14px;margin-right:6px" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
              </svg>
              Choose photo
            </span>
          </label>
          <span v-if="avatarFile" style="font-size:0.75rem;color:var(--text-3);margin-left:8px">{{ avatarFile.name }}</span>
        </div>
      </div>

      <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px">
        <div class="form-field">
          <label>First name</label>
          <input class="input" v-model="form.first_name" placeholder="Jane" />
        </div>
        <div class="form-field">
          <label>Last name</label>
          <input class="input" v-model="form.last_name" placeholder="Doe" />
        </div>
      </div>
      <div class="form-field">
        <label>Email</label>
        <input class="input" v-model="form.email" type="email" placeholder="you@example.com" />
      </div>
      <div class="form-field">
        <label>Password</label>
        <input class="input" v-model="form.password" type="password" placeholder="Min. 8 characters" />
      </div>
      <div class="form-field">
        <label>Date of birth</label>
        <input class="input" v-model="form.date_of_birth" type="date" />
      </div>
      <div class="form-field">
        <label>Nickname <span style="color:var(--text-3);font-weight:400">(optional)</span></label>
        <input class="input" v-model="form.nickname" placeholder="@handle" />
      </div>
      <div class="form-field">
        <label>About me <span style="color:var(--text-3);font-weight:400">(optional)</span></label>
        <textarea class="input" v-model="form.about_me" placeholder="Tell the world about yourself…" rows="2"></textarea>
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>

      <button class="btn btn-primary w-full" style="margin-top:20px" :disabled="loading" @click="register">
        {{ loading ? 'Creating account…' : 'Create account' }}
      </button>

      <p class="auth-footer">
        Already have an account? <router-link to="/login">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '../services/api'
import { setLoggedIn } from '../stores/auth'

const router = useRouter()
const error = ref('')
const loading = ref(false)
const avatarFile = ref(null)
const avatarPreview = ref('')
const avatarInput = ref(null)
const form = reactive({ email:'', password:'', first_name:'', last_name:'', date_of_birth:'', nickname:'', about_me:'' })

function pickAvatar(e) {
  const f = e.target.files?.[0]
  if (!f) return
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
  avatarFile.value = f
  avatarPreview.value = URL.createObjectURL(f)
}

async function register() {
  error.value = ''
  if (!form.email || !form.password || !form.first_name || !form.last_name || !form.date_of_birth) {
    error.value = 'Please fill in all required fields.'; return
  }
  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRe.test(form.email.trim())) {
    error.value = 'Please enter a valid email address.'; return
  }
  if (form.first_name.trim().length < 2 || form.first_name.trim().length > 50) {
    error.value = 'First name must be between 2 and 50 characters.'; return
  }
  if (form.last_name.trim().length < 2 || form.last_name.trim().length > 50) {
    error.value = 'Last name must be between 2 and 50 characters.'; return
  }
  if (form.password.length < 8) { error.value = 'Password must be at least 8 characters.'; return }
  if (form.password.length > 72) { error.value = 'Password cannot exceed 72 characters.'; return }
  const dob = new Date(form.date_of_birth)
  const minAge = new Date(); minAge.setFullYear(minAge.getFullYear() - 13)
  if (dob > minAge) { error.value = 'You must be at least 13 years old.'; return }
  if (dob > new Date()) { error.value = 'Date of birth cannot be in the future.'; return }
  if (form.nickname && form.nickname.trim().length > 30) {
    error.value = 'Nickname cannot exceed 30 characters.'; return
  }
  if (form.about_me && form.about_me.length > 500) {
    error.value = 'About me cannot exceed 500 characters.'; return
  }
  if (avatarFile.value) {
    const allowed = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
    if (!allowed.includes(avatarFile.value.type)) {
      error.value = 'Avatar must be a JPG, PNG, GIF or WebP image.'; return
    }
    if (avatarFile.value.size > 5 * 1024 * 1024) {
      error.value = 'Avatar file size cannot exceed 5 MB.'; return
    }
  }

  loading.value = true
  try {
    // Register first
    await apiFetch('/register', { method: 'POST', body: JSON.stringify(form) })

    // Then upload avatar if provided
    if (avatarFile.value) {
      const fd = new FormData()
      fd.append('avatar', avatarFile.value)
      await apiFetch('/upload/avatar', { method: 'POST', body: fd })
    }

    setLoggedIn(true)
    router.push('/feed')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onUnmounted(() => {
  if (avatarPreview.value) URL.revokeObjectURL(avatarPreview.value)
})
</script>
