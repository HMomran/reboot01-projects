import { reactive, ref } from 'vue'

/**
 * Shared chat state — owned by App.vue (which opens and maintains the WebSocket).
 * ChatView reads from this store and registers event handlers via onChatEvent().
 * Navbar reads unreadTotal for the badge.
 */
export const chatState = reactive({
  socket: null,
  connected: false,

  // Total unread messages across all conversations
  unreadTotal: 0,

  // Per-sender unread counts: { [senderID]: number }
  unreadsBySender: {},

  // Per-group unread counts: { [groupID]: number }
  unreadsByGroup: {},

  // The userID who is currently typing TO us (null if nobody)
  typingFrom: null,
  _typingTimer: null,
})

/**
 * The userID of the conversation currently open in ChatView.
 * App.vue passes this to handleGlobalChatMessage so messages from the
 * active conversation are not double-counted as unread.
 */
export const activeChatUserId = ref(null)

/**
 * The groupID of the group conversation currently open in ChatView.
 */
export const activeChatGroupId = ref(null)

// --- Event handler registry (used by ChatView) ---
const _handlers = []

/**
 * Register a callback that receives every raw WS event object.
 * Returns an unsubscribe function.
 */
export function onChatEvent(fn) {
  _handlers.push(fn)
  return () => {
    const idx = _handlers.indexOf(fn)
    if (idx !== -1) _handlers.splice(idx, 1)
  }
}

/** Internal: dispatch a WS event to all registered handlers. */
export function _dispatchChatEvent(event) {
  _handlers.forEach(fn => fn(event))
}

/**
 * Handle an incoming WS message at the global level:
 *   - type "message"  → increment unread count for sender
 *   - type "typing"   → set typingFrom + auto-clear after 3 s
 *   - type "read"     → (ChatView handles the seen-update via _dispatchChatEvent)
 *
 * @param {object} msg  parsed JSON from the WebSocket
 * @param {string|null} activeChatUserId  currently open chat user (skip unread increment)
 */
// currentUserId is set by ChatView so the store knows who we are
export let selfUserId = ''
export function setSelfUserId(id) { selfUserId = id }

export function handleGlobalChatMessage(msg, activeChatUserId) {
  if (msg.type === 'message') {
    // Never count our own sent messages as unread, but still dispatch
    // so ChatView can append the server-echoed message for the sender
    if (msg.sender_id && msg.sender_id === selfUserId) {
      _dispatchChatEvent(msg)
      return
    }

    // Only count as unread if not from the currently-open conversation
    if (msg.sender_id && msg.sender_id !== activeChatUserId) {
      chatState.unreadsBySender[msg.sender_id] =
        (chatState.unreadsBySender[msg.sender_id] || 0) + 1
      chatState.unreadTotal++
    }
  }

  if (msg.type === 'typing') {
    chatState.typingFrom = msg.sender_id
    clearTimeout(chatState._typingTimer)
    chatState._typingTimer = setTimeout(() => { chatState.typingFrom = null }, 3000)
  }

  if (msg.type === 'stop_typing') {
    if (chatState.typingFrom === msg.sender_id) {
      clearTimeout(chatState._typingTimer)
      chatState.typingFrom = null
    }
  }

  _dispatchChatEvent(msg)
}

/**
 * Clear unread count for a specific sender (call when opening that conversation).
 */
export function clearUnread(senderID) {
  const prev = chatState.unreadsBySender[senderID] || 0
  chatState.unreadTotal = Math.max(0, chatState.unreadTotal - prev)
  chatState.unreadsBySender[senderID] = 0
  delete chatState.unreadsBySender[senderID]
}

/**
 * Clear unread count for a specific group (call when opening that conversation).
 */
export function clearGroupUnread(groupID) {
  const prev = chatState.unreadsByGroup[groupID] || 0
  chatState.unreadTotal = Math.max(0, chatState.unreadTotal - prev)
  chatState.unreadsByGroup[groupID] = 0
  delete chatState.unreadsByGroup[groupID]
}
