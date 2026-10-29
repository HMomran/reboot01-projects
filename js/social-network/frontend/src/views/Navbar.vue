<template>
  <nav class="navbar">
    <div class="nav-left" style="display:flex;align-items:center;gap:12px">
      <!-- Outlined Navigation Icons -->
      <router-link to="/feed" class="nav-item" title="Home">
        <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
          <path class="outline-path" d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
          <polyline class="outline-path" points="9 22 9 12 15 12 15 22"/>
          <path class="fill-path" d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/>
        </svg>
      </router-link>

      <template v-if="isLoggedIn">
        <router-link to="/explore" class="nav-item" title="Explore">
          <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <circle class="outline-path" cx="12" cy="12" r="10"/>
            <polygon class="outline-path" points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"/>
            <path class="fill-path" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 11.5l-4.5 1.5 1.5-4.5L14.5 9 13 13.5z"/>
          </svg>
        </router-link>

        <router-link to="/chat" class="nav-item chat-nav-link" :class="{ shake: shaking }" @click="clearChatBadge" title="Chat">
          <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <path class="outline-path" d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            <path class="fill-path" d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
          </svg>
          <span v-if="chatState.unreadTotal > 0 && route.path !== '/chat'" class="chat-badge">
            {{ chatState.unreadTotal > 99 ? '99+' : chatState.unreadTotal }}
          </span>
        </router-link>

        <router-link to="/notifications" class="nav-item notif-nav-link" @click="clearNotifBadge" title="Notifications">
          <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <path class="outline-path" d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9zm-8 13a2 2 0 0 0 4 0"/>
            <path class="fill-path" d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.89 2 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.63 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
          </svg>
          <span v-if="notifState.unreadCount > 0" class="notif-badge">
            {{ notifState.unreadCount > 99 ? '99+' : notifState.unreadCount }}
          </span>
        </router-link>

        <router-link to="/groups" class="nav-item" title="Groups">
          <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <path class="outline-path" d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle class="outline-path" cx="9" cy="7" r="4"/>
            <path class="outline-path" d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path class="outline-path" d="M16 3.13a4 4 0 0 1 0 7.75"/>
            <path class="fill-path" d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
          </svg>
        </router-link>

        <router-link to="/profile" class="nav-item" title="Profile">
          <svg class="nav-icon" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <path class="outline-path" d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle class="outline-path" cx="12" cy="7" r="4"/>
            <path class="fill-path" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
          </svg>
        </router-link>
      </template>
    </div>

    <div class="nav-right" style="display:flex;align-items:center;gap:10px">
      <!-- Notification mute toggle -->
      <button v-if="isLoggedIn"
        class="icon-toggle"
        :title="notifMuted ? 'Notifications off — click to turn on' : 'Notifications on — click to mute'"
        @click="toggleNotifMute"
      >
        <svg v-if="notifMuted" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-bell-off"><path d="M13.73 21a2 2 0 0 1-3.46 0m10.19-2.28L4 3m14 5a6 6 0 0 0-9.33-5m2.09 13.73A6 6 0 0 1 6 8v5.58L3 19h15.28"></path></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-bell"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9zm-8 13a2 2 0 0 0 4 0"></path></svg>
      </button>

      <!-- Dark / Light theme toggle -->
      <button
        class="icon-toggle"
        :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
        @click="toggleTheme"
      >
        <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-sun"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-moon"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
      </button>

      <button v-if="isLoggedIn" @click="logout" class="nav-logout-btn">Logout</button>
      <template v-else>
        <router-link to="/login" class="btn btn-secondary btn-sm" style="font-size:0.8rem;padding:6px 12px;text-decoration:none">Log In</router-link>
        <router-link to="/register" class="btn btn-primary btn-sm" style="font-size:0.8rem;padding:6px 12px;text-decoration:none">Sign Up</router-link>
      </template>
    </div>
  </nav>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { apiFetch } from '../services/api'
import { chatState } from '../stores/chat'
import { notifState, resetNotifCount } from '../stores/notif'
import { isLoggedIn, setLoggedIn } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const shaking = ref(false)

// Shake chat icon briefly whenever new messages arrive
let prevTotal = chatState.unreadTotal
watch(() => chatState.unreadTotal, (val) => {
  if (val > prevTotal) {
    shaking.value = true
    setTimeout(() => { shaking.value = false }, 600)
  }
  prevTotal = val
})

function clearChatBadge() {}
function clearNotifBadge() { resetNotifCount() }

async function logout() {
  try { await apiFetch('/logout', { method: 'POST' }) } catch {}
  setLoggedIn(false)
  chatState.unreadTotal = 0
  chatState.unreadsBySender = {}
  chatState.connected = false
  notifState.unreadCount = 0
  router.push('/login')
}

// ── Notification mute ─────────────────────────────────────
const notifMuted = ref(localStorage.getItem('notif_muted') === 'true')
function toggleNotifMute() {
  notifMuted.value = !notifMuted.value
  localStorage.setItem('notif_muted', notifMuted.value)
  window.dispatchEvent(new CustomEvent('notif-mute-changed', { detail: notifMuted.value }))
}

// ── Theme toggle ──────────────────────────────────────────
const isDark = ref(true)

function applyTheme(dark) {
  if (dark) {
    document.documentElement.removeAttribute('data-theme')
  } else {
    document.documentElement.setAttribute('data-theme', 'light')
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  applyTheme(isDark.value)
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  isDark.value = saved !== 'light'
  applyTheme(isDark.value)
})
</script>

<style scoped>
.nav-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  color: var(--text-2);
  transition: all 0.2s ease;
  position: relative;
}

.nav-item:hover {
  color: var(--text);
  background: var(--surface2);
}

.nav-item.router-link-active {
  color: var(--accent);
  background: var(--accent-bg);
}

.nav-icon {
  width: 22px;
  height: 22px;
  stroke: currentColor;
  fill: none;
  stroke-width: 2px;
  transition: transform 0.2s ease;
}

.nav-item:hover .nav-icon {
  transform: scale(1.06);
}

/* outline elements visible by default */
.nav-icon .outline-path {
  display: block;
}

/* fill elements hidden by default */
.nav-icon .fill-path {
  display: none;
}

/* active link swaps outline for fill */
.nav-item.router-link-active .nav-icon .outline-path {
  display: none;
}

.nav-item.router-link-active .nav-icon .fill-path {
  display: block;
  fill: currentColor;
}

.chat-nav-link {
  position: relative;
}

.chat-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: var(--danger);
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
  border-radius: 10px;
  padding: 0px 4px;
  min-width: 14px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  animation: badge-pop 0.25s ease;
  pointer-events: none;
}

@keyframes badge-pop {
  0%   { transform: scale(0.6); }
  70%  { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.notif-nav-link {
  position: relative;
}

.notif-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: var(--danger);
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
  border-radius: 10px;
  padding: 0px 4px;
  min-width: 14px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  animation: badge-pop 0.25s ease;
  pointer-events: none;
}

.shake {
  animation: nav-shake 0.5s ease;
}

@keyframes nav-shake {
  0%, 100% { transform: translateX(0); }
  20%       { transform: translateX(-3px); }
  40%       { transform: translateX(3px); }
  60%       { transform: translateX(-2px); }
  80%       { transform: translateX(2px); }
}

.icon-toggle {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.15rem;
  padding: 6px;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  transition: background 0.15s, color 0.15s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-2);
}
.icon-toggle:hover { background: var(--surface2); color: var(--text); }

.nav-logout-btn {
  padding: 7px 16px;
  border-radius: var(--radius-sm);
  background: var(--surface2);
  border: 1px solid var(--border2);
  color: var(--text-2);
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s;
}
.nav-logout-btn:hover { color: var(--danger); border-color: var(--danger); background: rgba(224,64,64,0.08); }
</style>
