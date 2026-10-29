<template>
  <router-view />

  <!-- Global toast notifications — visible on every page -->
  <div style="position:fixed;bottom:24px;right:24px;display:flex;flex-direction:column;gap:10px;z-index:1000;pointer-events:none">
    <div v-for="toast in toasts" :key="toast.id"
      style="background:var(--surface);border:1px solid var(--border2);border-radius:12px;padding:14px 18px;box-shadow:var(--shadow-lg);max-width:320px;display:flex;align-items:flex-start;gap:12px;animation:slideIn 0.3s ease;pointer-events:all">
      <span style="font-size:1.2rem;flex-shrink:0">{{ toastIcon(toast.type) }}</span>
      <p style="font-size:0.85rem;color:var(--text);margin:0;flex:1;line-height:1.5">{{ toast.message }}</p>
      <button @click="removeToast(toast.id)" style="background:none;border:none;cursor:pointer;color:var(--text-3);font-size:0.9rem;padding:0;flex-shrink:0">✕</button>
    </div>
  </div>
</template>

<style>
@keyframes slideIn {
  from { transform: translateX(110px); opacity: 0; }
  to   { transform: translateX(0);     opacity: 1; }
}
</style>

<script setup>
import { onMounted, onUnmounted, watch, ref } from 'vue'
import { useRoute } from 'vue-router'
import { apiFetch } from './services/api'
import { chatState, handleGlobalChatMessage, activeChatUserId, selfUserId } from './stores/chat'
import { notifState, _dispatchNotifEvent, setNotifCount, onNotifEvent, addToast, toasts, removeToast } from './stores/notif'
import { isLoggedIn } from './stores/auth'

const route = useRoute()
let presenceSocket = null
let presenceRetry = null
let notifSocket = null
let notifRetry = null
let _destroyed = false

// Notification mute toggle — persisted in localStorage
const notifMuted = ref(localStorage.getItem('notif_muted') === 'true')
function toggleNotifMute() {
  notifMuted.value = !notifMuted.value
  localStorage.setItem('notif_muted', notifMuted.value)
}

const PUBLIC_PATHS = ['/login', '/register']

function isPublic() { return PUBLIC_PATHS.includes(route.path) || !isLoggedIn.value }

function connectPresence() {
  if (_destroyed || isPublic()) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  presenceSocket = new WebSocket(`${proto}//${location.host}/ws/presence`)
  presenceSocket.onclose = () => {
    if (!_destroyed && !isPublic()) presenceRetry = setTimeout(connectPresence, 5000)
  }
}

function connectChatWS() {
  if (_destroyed || isPublic()) return
  if (chatState.socket && chatState.socket.readyState < 2) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const ws = new WebSocket(`${proto}//${location.host}/ws/private`)
  chatState.socket = ws
  ws.onopen = () => { chatState.connected = true }
  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      handleGlobalChatMessage(msg, activeChatUserId.value)
      // Live toast for incoming messages (only if NOT on the chat page and not muted)
      if (msg.type === 'message' && msg.sender_id !== selfUserId
          && route.path !== '/chat' && !notifMuted.value) {
        const sender = msg.sender_name || 'Someone'
        const text = msg.content && msg.content.length > 60 ? msg.content.substring(0, 60) + '…' : msg.content
        addToast({ type: 'message', message: sender + ': ' + text })
      }
    } catch {}
  }
  ws.onclose = () => {
    chatState.connected = false
    chatState.socket = null
    if (!_destroyed && !isPublic()) setTimeout(connectChatWS, 5000)
  }
}

function connectNotifWS() {
  if (_destroyed || isPublic()) return
  if (notifSocket && notifSocket.readyState < 2) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const ws = new WebSocket(`${proto}//${location.host}/ws/notifications`)
  notifState.socket = ws
  ws.onopen = () => { notifState.connected = true }
  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (msg.event === 'notification') _dispatchNotifEvent(msg)
    } catch {}
  }
  ws.onclose = () => {
    notifState.connected = false
    notifState.socket = null
    notifSocket = null
    if (!_destroyed && !isPublic()) notifRetry = setTimeout(connectNotifWS, 5000)
  }
  notifSocket = ws
}

watch(() => [route.path, isLoggedIn.value], ([path, loggedIn]) => {
  if (!PUBLIC_PATHS.includes(path) && loggedIn) {
    connectPresence()
    connectChatWS()
    connectNotifWS()
  } else {
    if (presenceSocket) { presenceSocket.onclose=null; presenceSocket.close(); presenceSocket=null }
    if (chatState.socket) { chatState.socket.onclose=null; chatState.socket.close(); chatState.socket=null; chatState.connected=false }
    if (notifSocket) { notifSocket.onclose=null; notifSocket.close(); notifSocket=null; notifState.connected=false; notifState.socket=null }
  }
}, { immediate: true })

const NOTIF_ICONS = { follow_request:'👤', new_follower:'✅', group_invite:'👥', group_join_request:'📨', group_event_created:'📅', like:'❤️', comment:'💬', message:'✉️', group_message:'👥', follow_accept:'✅', group_accept:'🎉' }
function toastIcon(t) { return NOTIF_ICONS[t] || '🔔' }

let _unsubNotif = null
onMounted(async () => {
  if (isLoggedIn.value) {
    try {
      const data = await apiFetch('/notifications/unread')
      setNotifCount(typeof data === 'number' ? data : data?.unread ?? 0)
    } catch {}
  }
  // Show toasts for live notifications on every page except the notifications page itself
  _unsubNotif = onNotifEvent((event) => {
    if (route.path !== '/notifications' && !notifMuted.value) addToast(event)
  })
})

onUnmounted(() => {
  _destroyed = true
  clearTimeout(presenceRetry)
  clearTimeout(notifRetry)
  if (presenceSocket) presenceSocket.close(1000)
  if (chatState.socket) chatState.socket.close(1000)
  if (notifSocket) notifSocket.close(1000)
  if (_unsubNotif) _unsubNotif()
})
</script>
