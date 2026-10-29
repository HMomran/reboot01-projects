import { reactive, ref } from 'vue'

export const notifState = reactive({
  socket: null,
  connected: false,
  unreadCount: 0,
})

// ── Global toast queue (shown in App.vue, visible on every page) ──────────────
export const toasts = ref([])
let _toastId = 0

export function addToast(notif) {
  const id = ++_toastId
  toasts.value.push({ id, type: notif.type || notif.event, message: notif.message })
  setTimeout(() => removeToast(id), 5000)
}

export function removeToast(id) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}

// ── Event handler registry ────────────────────────────────────────────────────
const _handlers = []
const _seenIds = new Set()

export function onNotifEvent(fn) {
  _handlers.push(fn)
  return () => {
    const i = _handlers.indexOf(fn)
    if (i !== -1) _handlers.splice(i, 1)
  }
}

import { chatState, activeChatGroupId } from './chat'

export function _dispatchNotifEvent(event) {
  // Deduplicate — WS can occasionally re-deliver the same event
  if (event.id && _seenIds.has(event.id)) return
  if (event.id) _seenIds.add(event.id)

  if (event.type === 'group_message') {
    const gId = event.group_id
    if (gId && gId !== activeChatGroupId.value) {
      if (!chatState.unreadsByGroup) chatState.unreadsByGroup = {}
      chatState.unreadsByGroup[gId] = (chatState.unreadsByGroup[gId] || 0) + 1
      chatState.unreadTotal++
    }
  } else {
    notifState.unreadCount = Math.max(0, notifState.unreadCount + 1)
  }
  _handlers.forEach(fn => fn(event))
}

export function resetNotifCount() {
  notifState.unreadCount = 0
}

export function setNotifCount(n) {
  notifState.unreadCount = typeof n === 'number' ? n : 0
}
