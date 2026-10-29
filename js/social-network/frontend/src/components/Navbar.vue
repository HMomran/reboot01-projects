<template>
  <nav class="navbar">
    <div class="nav-left">
      <router-link to="/feed">Home</router-link>
      <router-link to="/explore">Explore</router-link>
      <router-link to="/chat" class="chat-nav-link" :class="{ shake: shaking }" @click="clearChatBadge">
        Chat
        <span v-if="chatState.unreadTotal > 0" class="chat-badge">
          {{ chatState.unreadTotal > 99 ? '99+' : chatState.unreadTotal }}
        </span>
      </router-link>
      <router-link to="/notifications" class="notif-nav-link" @click="clearNotifBadge">
        Notifications
        <span v-if="notifState.unreadCount > 0" class="notif-badge">
          {{ notifState.unreadCount > 99 ? '99+' : notifState.unreadCount }}
        </span>
      </router-link>
      <router-link to="/groups">Groups</router-link>
    </div>

    <div class="nav-right">
      <span class="nav-avatar" @click="goToProfile" title="Profile">👤</span>
      <button @click="logout">Logout</button>
    </div>
  </nav>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '../services/api'
import { chatState, clearUnread } from '../stores/chat'
import { notifState, resetNotifCount } from '../stores/notif'
import { setLoggedIn } from '../stores/auth'

const router = useRouter()
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

function clearChatBadge() {
  // Counts are cleared per-conversation when the user opens one in ChatView.
  // Clicking the nav link only visually resets the top-level badge if all
  // per-sender counts have already been cleared by ChatView.
}

function clearNotifBadge() {
  resetNotifCount()
}

function goToProfile() {
  router.push('/profile')
}

async function logout() {
  try {
    await apiFetch('/logout', { method: 'POST' })
  } catch {}
  setLoggedIn(false)
  chatState.unreadTotal = 0
  chatState.unreadsBySender = {}
  chatState.connected = false
  notifState.unreadCount = 0
  router.push('/login')
}
</script>

<style scoped>
.chat-nav-link {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.chat-badge {
  position: absolute;
  top: -8px;
  right: -14px;
  background: #ef4444;
  color: #fff;
  font-size: 0.65rem;
  font-weight: 700;
  border-radius: 10px;
  padding: 1px 5px;
  min-width: 16px;
  text-align: center;
  line-height: 1.4;
  animation: badge-pop 0.25s ease;
}

@keyframes badge-pop {
  0%   { transform: scale(0.6); }
  70%  { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.notif-nav-link {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.notif-badge {
  position: absolute;
  top: -8px;
  right: -14px;
  background: #ef4444;
  color: #fff;
  font-size: 0.65rem;
  font-weight: 700;
  border-radius: 10px;
  padding: 1px 5px;
  min-width: 16px;
  text-align: center;
  line-height: 1.4;
  animation: badge-pop 0.25s ease;
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
</style>


