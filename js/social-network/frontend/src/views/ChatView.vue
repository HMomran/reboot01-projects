<template>
  <div>
    <Navbar />
    <div class="page-layout two-col" style="height:calc(100vh - var(--nav-h));padding-top:20px;padding-bottom:20px">

      <!-- Sidebar -->
      <div class="card" style="display:flex;flex-direction:column;overflow:hidden;height:100%">

        <!-- Tabs: DMs / Groups -->
        <div style="display:flex;border-bottom:1px solid var(--border);flex-shrink:0;align-items:center;justify-content:space-between">
          <div style="display:flex;flex:1">
            <button @click="tab='dm'" :style="{ flex:1, padding:'12px', background:'none', border:'none', borderBottom: tab==='dm' ? '2px solid var(--accent)' : '2px solid transparent', color: tab==='dm' ? 'var(--accent)' : 'var(--text-2)', fontWeight:'600', fontSize:'0.85rem', cursor:'pointer', fontFamily:'inherit', display:'inline-flex', alignItems:'center', justifyContent:'center', gap:'4px' }">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
              Direct
              <span v-if="tab !== 'dm' && hasUnreadDms" style="width:6px;height:6px;background:var(--accent);border-radius:50%;display:inline-block;margin-left:4px"></span>
            </button>
            <button @click="tab='group'" :style="{ flex:1, padding:'12px', background:'none', border:'none', borderBottom: tab==='group' ? '2px solid var(--accent)' : '2px solid transparent', color: tab==='group' ? 'var(--accent)' : 'var(--text-2)', fontWeight:'600', fontSize:'0.85rem', cursor:'pointer', fontFamily:'inherit', display:'inline-flex', alignItems:'center', justifyContent:'center', gap:'4px' }">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
              Groups
              <span v-if="tab !== 'group' && hasUnreadGroups" style="width:6px;height:6px;background:var(--accent);border-radius:50%;display:inline-block;margin-left:4px"></span>
            </button>
          </div>
          <button v-if="tab==='dm'" @click="showNewChatModal = true" class="btn btn-ghost" style="padding:12px;border:none;background:none;color:var(--accent);cursor:pointer;display:inline-flex;align-items:center" title="New Chat">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4z"></path></svg>
          </button>
        </div>

        <div style="padding:10px;border-bottom:1px solid var(--border);flex-shrink:0">
          <input class="input" v-model="search" :placeholder="tab==='dm' ? 'Search chats…' : 'Search groups…'" style="font-size:0.85rem" />
        </div>

        <div style="flex:1;overflow-y:auto;padding:6px">

          <!-- DM tab -->
          <template v-if="tab === 'dm'">
            <template v-if="filteredRecentConvos.length">
              <div v-for="c in filteredRecentConvos" :key="c.user.id"
                :class="['sidebar-item', activeType==='dm' && activeDmUserId===c.user.id && 'active']"
                @click="openDm(c.user)">
                <div style="position:relative;flex-shrink:0">
                  <div class="avatar-ring" style="width:36px;height:36px;font-size:0.75rem" :style="{ background: avatarColor(c.user.id) }">
                    <img v-if="c.user.avatar" :src="BACKEND_URL + c.user.avatar" />
                    <span v-else>{{ initials(c.user.first_name,c.user.last_name) }}</span>
                  </div>
                  <span v-if="onlineIds.has(c.user.id)" class="presence-dot"></span>
                </div>
                <div style="flex:1;overflow:hidden">
                  <div style="display:flex;align-items:center;justify-content:space-between">
                    <span style="font-size:0.87rem;font-weight:600;color:var(--text)">{{ c.user.first_name }} {{ c.user.last_name }}</span>
                    <span v-if="chatState.unreadsBySender[c.user.id] > 0" class="badge">{{ chatState.unreadsBySender[c.user.id] }}</span>
                  </div>
                  <span style="font-size:0.74rem;color:var(--text-3);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;display:block">{{ c.lastMessage || 'No messages yet' }}</span>
                </div>
              </div>
            </template>
            <div v-else style="color:var(--text-3);font-size:0.82rem;padding:12px 8px;text-align:center">
              No chats yet.<br>Click the edit icon above to start a chat!
            </div>
          </template>

          <!-- Group tab -->
          <template v-if="tab === 'group'">
            <div v-if="!filteredGroups.length" style="color:var(--text-3);font-size:0.82rem;padding:12px 8px">You haven't joined any groups yet.</div>
            <div v-for="g in filteredGroups" :key="g.id"
              :class="['sidebar-item', activeType==='group' && activeGroupId===g.id && 'active']"
              @click="openGroup(g)">
              <div style="width:36px;height:36px;border-radius:10px;background:var(--accent-bg);display:flex;align-items:center;justify-content:center;color:var(--accent);flex-shrink:0">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
              </div>
              <div style="flex:1;overflow:hidden">
                <div style="display:flex;align-items:center;justify-content:space-between">
                  <div style="font-size:0.87rem;font-weight:600;color:var(--text);white-space:nowrap;overflow:hidden;text-overflow:ellipsis">{{ g.title }}</div>
                  <span v-if="chatState.unreadsByGroup[g.id] > 0" class="badge">{{ chatState.unreadsByGroup[g.id] }}</span>
                </div>
                <div style="font-size:0.74rem;color:var(--text-3)">{{ g.member_count }} members</div>
              </div>
            </div>
          </template>
        </div>
      </div>

      <!-- Chat window -->
      <div class="card" style="display:flex;flex-direction:column;overflow:hidden;height:100%" v-if="activeType">

        <!-- DM Header -->
        <div v-if="activeType==='dm'" style="padding:14px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:12px;flex-shrink:0">
          <div class="avatar-ring" style="width:38px;height:38px;font-size:0.78rem;flex-shrink:0" :style="{ background: avatarColor(activeDmUserId) }">
            <img v-if="activeUser?.avatar" :src="BACKEND_URL + activeUser.avatar" />
            <span v-else>{{ activeUser ? initials(activeUser.first_name, activeUser.last_name) : '?' }}</span>
          </div>
          <div>
            <div style="font-weight:700;font-size:0.95rem">{{ activeUser ? activeUser.first_name + ' ' + activeUser.last_name : '' }}</div>
            <div style="font-size:0.74rem;display:flex;align-items:center;gap:5px">
              <template v-if="typingFrom===activeDmUserId">
                <span class="typing-dots"><span></span><span></span><span></span></span>
                <em style="color:var(--accent);font-style:normal">typing…</em>
              </template>
              <template v-else-if="onlineIds.has(activeDmUserId)">
                <span class="status-online"></span><span style="color:var(--success)">Online</span>
              </template>
              <template v-else>
                <span class="status-offline"></span><span style="color:var(--text-3)">Offline</span>
              </template>
            </div>
          </div>
        </div>

        <!-- Group Header -->
        <div v-if="activeType==='group'" style="padding:14px 18px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:12px;flex-shrink:0">
          <div style="width:38px;height:38px;border-radius:10px;background:var(--accent-bg);display:flex;align-items:center;justify-content:center;color:var(--accent);flex-shrink:0">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
          </div>
          <div>
            <div style="font-weight:700;font-size:0.95rem">{{ activeGroupTitle }}</div>
            <div style="font-size:0.74rem;color:var(--text-3)">Group chat</div>
          </div>
        </div>

        <!-- Messages -->
        <div ref="chatBox" @scroll="onChatScroll" style="flex:1;overflow-y:auto;padding:18px;display:flex;flex-direction:column;gap:6px">
          <!-- Load older indicator -->
          <div v-if="loadingOlder" style="text-align:center;padding:8px;font-size:0.78rem;color:var(--text-3)">Loading older messages…</div>
          <div v-else-if="noMoreMsgs && messages.length > 0" style="text-align:center;padding:8px;font-size:0.75rem;color:var(--text-3)">Beginning of conversation</div>
          <div v-if="loadingMsgs" v-for="n in 4" :key="n" style="display:flex" :style="{ justifyContent: n%2===0 ? 'flex-end' : 'flex-start' }">
            <div class="skeleton" style="height:36px;border-radius:12px;width:48%"></div>
          </div>
          <template v-else>
            <div v-for="msg in messages" :key="msg.id"
              style="display:flex;flex-direction:column;max-width:72%"
              :style="{ alignSelf: msg.sender_id===currentUserId ? 'flex-end' : 'flex-start', alignItems: msg.sender_id===currentUserId ? 'flex-end' : 'flex-start' }">
              <!-- Group: show sender name -->
              <div v-if="activeType==='group' && msg.sender_id !== currentUserId" style="font-size:0.72rem;color:var(--text-3);margin-bottom:2px;padding-left:4px">{{ getSenderName(msg.sender_id) }}</div>
              <div :style="{ background: msg.sender_id===currentUserId ? 'var(--accent)' : 'var(--surface2)', color: msg.sender_id===currentUserId ? '#000' : 'var(--text)', borderRadius: msg.sender_id===currentUserId ? '14px 14px 4px 14px' : '14px 14px 14px 4px', padding:'9px 14px', fontSize:'0.91rem', lineHeight:'1.5', wordBreak:'break-word' }">
                {{ msg.content }}
              </div>
              <div style="display:flex;align-items:center;gap:4px;margin-top:3px">
                <span style="font-size:0.68rem;color:var(--text-3)">{{ relTime(msg.created_at) }}</span>
                <span v-if="activeType==='dm' && msg.sender_id===currentUserId" style="font-size:0.7rem;color:var(--accent);font-weight:700">{{ msg.is_read ? '✓✓' : '✓' }}</span>
              </div>
            </div>
            <div v-if="activeType==='dm' && typingFrom===activeDmUserId" style="align-self:flex-start">
              <div style="background:var(--surface2);border-radius:14px 14px 14px 4px;padding:10px 16px">
                <span class="typing-dots"><span></span><span></span><span></span></span>
              </div>
            </div>
          </template>
        </div>

        <!-- Compose -->
        <div style="padding:12px 16px;border-top:1px solid var(--border);flex-shrink:0">
          <!-- Emoji picker -->
          <div v-if="showEmoji" style="background:var(--surface2);border:1px solid var(--border);border-radius:12px;padding:10px;margin-bottom:8px;display:flex;flex-wrap:wrap;gap:6px;max-height:140px;overflow-y:auto">
            <button v-for="e in EMOJIS" :key="e" @click="insertEmoji(e)"
              style="background:none;border:none;cursor:pointer;font-size:1.3rem;padding:2px 4px;border-radius:6px;line-height:1;transition:background 0.1s"
              @mouseenter="ev=>ev.target.style.background='var(--surface3)'"
              @mouseleave="ev=>ev.target.style.background='none'">{{ e }}</button>
          </div>
          <div style="display:flex;gap:8px;align-items:center">
            <button @click="showEmoji=!showEmoji"
              style="background:none;border:none;cursor:pointer;padding:6px;border-radius:8px;color:var(--text-2);transition:background 0.15s;flex-shrink:0;display:inline-flex;align-items:center;justify-content:center"
              :style="{ background: showEmoji ? 'var(--surface2)' : 'none' }">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M8 14s1.5 2 4 2 4-2 4-2"></path><line x1="9" y1="9" x2="9.01" y2="9"></line><line x1="15" y1="9" x2="15.01" y2="9"></line></svg>
            </button>
            <div style="flex:1;position:relative">
              <input class="input" v-model="newMsg" ref="msgInput" placeholder="Type a message…" style="width:100%;border-radius:22px;box-sizing:border-box" @keyup.enter="sendMsg" @input="onTyping"
                :style="{ borderColor: newMsg.length > 1000 ? 'var(--danger)' : '' }" />
              <span v-if="newMsg.length > 900" style="position:absolute;right:12px;top:50%;transform:translateY(-50%);font-size:0.68rem;pointer-events:none"
                :style="{ color: newMsg.length > 1000 ? 'var(--danger)' : 'var(--text-3)' }">{{ newMsg.length }}/1000</span>
            </div>
            <button class="btn btn-primary" style="border-radius:22px;padding:10px;width:38px;height:38px;display:inline-flex;align-items:center;justify-content:center;flex-shrink:0" :disabled="!newMsg.trim() || newMsg.length > 1000" @click="sendMsg">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="19" x2="12" y2="5"></line><polyline points="5 12 12 5 19 12"></polyline></svg>
            </button>
          </div>
        </div>
        <p v-if="chatError" style="color:var(--danger);font-size:0.8rem;padding:0 16px 8px">{{ chatError }}</p>
      </div>

      <!-- Empty state -->
      <div class="card empty-state" style="height:100%" v-else>
        <div class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
        </div>
        <h3>Select a conversation</h3>
        <p>Choose a person or group to start chatting</p>
      </div>
    </div>

    <!-- New Chat Modal -->
    <div v-if="showNewChatModal" class="modal-backdrop" @click.self="showNewChatModal=false">
      <div class="modal-box" style="max-width:400px">
        <div class="modal-header">
          <h3>New Message</h3>
          <button class="modal-close" @click="showNewChatModal=false">✕</button>
        </div>
        <div style="padding:10px;border-bottom:1px solid var(--border)">
          <input class="input" v-model="modalSearch" placeholder="Search your connections…" style="font-size:0.85rem;width:100%" />
        </div>
        <div style="max-height:300px;overflow-y:auto;padding:8px">
          <div v-if="!modalFilteredUsers.length" style="color:var(--text-3);text-align:center;padding:24px;font-size:0.88rem">
            No contacts found.
          </div>
          <div v-for="u in modalFilteredUsers" :key="u.id" 
            style="display:flex;align-items:center;gap:10px;padding:8px;cursor:pointer;border-radius:8px;transition:background 0.14s"
            @click="startChatWith(u)"
            @mouseenter="ev=>ev.currentTarget.style.background='var(--surface2)'"
            @mouseleave="ev=>ev.currentTarget.style.background='none'">
            <div class="avatar-ring" style="width:36px;height:36px;font-size:0.75rem;flex-shrink:0" :style="{ background: avatarColor(u.id) }">
              <img v-if="u.avatar" :src="BACKEND_URL + u.avatar" />
              <span v-else>{{ initials(u.first_name,u.last_name) }}</span>
            </div>
            <div style="flex:1">
              <span style="font-weight:600;font-size:0.9rem;color:var(--text)">{{ u.first_name }} {{ u.last_name }}</span>
              <div v-if="u.nickname" style="font-size:0.72rem;color:var(--text-3)">@{{ u.nickname }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'
import { chatState, onChatEvent, clearUnread, clearGroupUnread, setSelfUserId, activeChatUserId as globalActiveChatUserId, activeChatGroupId as globalActiveChatGroupId } from '../stores/chat'

const router = useRouter()
const route = useRoute()

// Tab state
const tab = ref('dm')
const activeType = ref(null) // 'dm' | 'group'
const activeDmUserId = ref(null)
const activeGroupId = ref(null)
const activeGroupTitle = ref('')

// Modal state
const showNewChatModal = ref(false)
const modalSearch = ref('')
const followingUsers = ref([])

// Data
const users = ref([])
const myGroups = ref([])
const recentConvos = ref([])
const search = ref('')
const messages = ref([])
const newMsg = ref('')
const chatError = ref('')
const currentUserId = ref('')
const loadingMsgs = ref(false)
const loadingOlder = ref(false)
const noMoreMsgs = ref(false)
const chatBox = ref(null)
let _scrollCooldown = false
const msgInput = ref(null)
const onlineIds = ref(new Set())
const showEmoji = ref(false)

// Group WS
let groupSocket = null

const EMOJIS = ['😀','😂','😍','🥰','😎','🤔','😢','😡','👍','👎','❤️','🔥','🎉','💯','✅','❌','🙏','💪','😭','🤣','😊','🥳','💀','👀','🫡','😏','🤝','👋','🙌','💬','📸','🎶','⚡','🌟','💫','🍕','🎮','🚀','💡','🎯']

const unreadsBySender = chatState.unreadsBySender
const typingFrom = computed(() => chatState.typingFrom)

let typingTimeout = null
function sendStopTyping() {
  if (activeType.value !== 'dm' || !activeDmUserId.value) return
  if (chatState.socket?.readyState === WebSocket.OPEN)
    chatState.socket.send(JSON.stringify({ type:'stop_typing', receiver_id:activeDmUserId.value }))
  clearTimeout(typingTimeout); typingTimeout = null
}
function onTyping() {
  if (activeType.value !== 'dm' || !activeDmUserId.value || !chatState.socket) return
  if (!newMsg.value) { sendStopTyping(); return }
  if (typingTimeout) return
  chatState.socket.send(JSON.stringify({ type:'typing', receiver_id:activeDmUserId.value }))
  typingTimeout = setTimeout(() => { typingTimeout = null }, 2000)
}

function insertEmoji(e) {
  newMsg.value += e
  showEmoji.value = false
  nextTick(() => msgInput.value?.focus())
}

const activeUser = computed(() => users.value.find(u => u.id === activeDmUserId.value) || null)

const hasUnreadDms = computed(() => {
  return Object.values(chatState.unreadsBySender).some(count => count > 0)
})

const hasUnreadGroups = computed(() => {
  return chatState.unreadsByGroup && Object.values(chatState.unreadsByGroup).some(count => count > 0)
})

const sortedRecentConvos = computed(() => {
  return [...recentConvos.value].sort((a, b) => {
    const timeA = a.lastAt ? new Date(a.lastAt.includes('T') ? a.lastAt : a.lastAt.replace(' ','T')+'Z').getTime() : 0
    const timeB = b.lastAt ? new Date(b.lastAt.includes('T') ? b.lastAt : b.lastAt.replace(' ','T')+'Z').getTime() : 0
    return timeB - timeA
  })
})

const filteredRecentConvos = computed(() => {
  return sortedRecentConvos.value.filter(c =>
    `${c.user.first_name} ${c.user.last_name} ${c.user.nickname || ''}`.toLowerCase().includes(search.value.toLowerCase())
  )
})

const modalFilteredUsers = computed(() => {
  const recentIds = new Set(recentConvos.value.map(c => c.user.id))
  return followingUsers.value.filter(u =>
    !recentIds.has(u.id) &&
    `${u.first_name} ${u.last_name} ${u.nickname || ''}`.toLowerCase().includes(modalSearch.value.toLowerCase())
  )
})

const filteredGroups = computed(() => myGroups.value.filter(g => g.title.toLowerCase().includes(search.value.toLowerCase())))

const COLORS = ['#f5a623','#3ecf78','#3b8aff','#e04040','#a78bfa','#f472b6']
function avatarColor(id) { let n=0; for(let c of (id||'')) n+=c.charCodeAt(0); return COLORS[n%COLORS.length] }
function initials(f,l) { return `${(f||'?')[0]}${(l||'?')[0]}`.toUpperCase() }
function relTime(s) {
  if (!s) return ''
  const normalized = s.includes('T') ? s : s.replace(' ','T')+'Z'
  const d = Date.now() - new Date(normalized).getTime()
  if (isNaN(d)) return ''
  const sec=Math.floor(d/1000); if(sec<60) return 'just now'
  const m=Math.floor(sec/60); if(m<60) return m+'m ago'
  const h=Math.floor(m/60); if(h<24) return h+'h ago'
  return Math.floor(h/24)+'d ago'
}

function getSenderName(id) {
  if (id === currentUserId.value) return 'You'
  const u = users.value.find(u => u.id === id)
  return u ? `${u.first_name} ${u.last_name}` : id.substring(0,8)
}

// ── DM ──────────────────────────────────────────────────
async function openDm(user) {
  sendStopTyping()
  closeGroupSocket()
  chatError.value = ''; showEmoji.value = false
  activeType.value = 'dm'; activeDmUserId.value = user.id
  globalActiveChatUserId.value = user.id
  globalActiveChatGroupId.value = null
  clearUnread(user.id)
  apiFetch('/messages/mark-read', { method:'POST', body:JSON.stringify({ sender_id:user.id }) }).catch(()=>{})
  if (chatState.socket?.readyState === WebSocket.OPEN)
    chatState.socket.send(JSON.stringify({ type:'read', sender_id:user.id }))
  loadingMsgs.value = true; messages.value = []; noMoreMsgs.value = false
  try {
    const d = await apiFetch(`/messages/private?user_id=${user.id}`)
    messages.value = Array.isArray(d) ? d : []
    // Always show this conversation at the top of the sidebar with the latest message
    const lastMsg = messages.value[messages.value.length - 1]
    const idx = recentConvos.value.findIndex(c => c.user.id === user.id)
    if (idx !== -1) {
      const c = recentConvos.value.splice(idx, 1)[0]
      c.lastMessage = lastMsg ? lastMsg.content : ''
      c.lastAt = lastMsg ? lastMsg.created_at : new Date().toISOString()
      recentConvos.value.unshift(c)
    } else {
      recentConvos.value.unshift({
        user,
        lastMessage: lastMsg ? lastMsg.content : '',
        lastAt: lastMsg ? lastMsg.created_at : new Date().toISOString()
      })
    }
  } catch(e) { chatError.value = e.message }
  finally { loadingMsgs.value = false; await scrollBottom() }
}

// ── Group chat ───────────────────────────────────────────
async function openGroup(group) {
  sendStopTyping()
  closeGroupSocket()
  chatError.value = ''; showEmoji.value = false
  activeType.value = 'group'; activeGroupId.value = group.id; activeGroupTitle.value = group.title
  globalActiveChatUserId.value = null
  globalActiveChatGroupId.value = group.id
  clearGroupUnread(group.id)
  loadingMsgs.value = true; messages.value = []; noMoreMsgs.value = false
  try { const d = await apiFetch(`/messages/group?group_id=${group.id}`); messages.value = Array.isArray(d)?d:[] }
  catch(e) { chatError.value = e.message } finally { loadingMsgs.value = false; await scrollBottom() }

  connectGroupSocket(group.id)
}

let _currentGroupId = null

function connectGroupSocket(groupId) {
  _currentGroupId = groupId
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const ws = new WebSocket(`${proto}//${location.host}/ws/group?group_id=${groupId}`)
  groupSocket = ws

  ws.onopen = () => { chatError.value = '' }

  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      // Avoid duplicate — backend echoes to all including sender
      const exists = msg.id && messages.value.some(m => m.id === msg.id)
      if (!exists) { messages.value.push(msg); scrollBottom() }
    } catch {}
  }

  ws.onerror = () => { chatError.value = 'Group connection lost — retrying…' }

  ws.onclose = () => {
    groupSocket = null
    // Auto-reconnect only if still viewing same group
    if (_currentGroupId === groupId) {
      setTimeout(() => {
        if (_currentGroupId === groupId) connectGroupSocket(groupId)
      }, 3000)
    }
  }
}

function closeGroupSocket() {
  _currentGroupId = null
  globalActiveChatGroupId.value = null
  if (groupSocket) {
    groupSocket.onclose = null  // prevent auto-reconnect
    groupSocket.close()
    groupSocket = null
  }
}

// ── Send ─────────────────────────────────────────────────
function sendMsg() {
  const content = newMsg.value.trim(); if (!content || newMsg.value.length > 1000) return
  showEmoji.value = false

  if (activeType.value === 'dm') {
    if (!chatState.socket || chatState.socket.readyState !== WebSocket.OPEN) {
      chatError.value = 'Not connected — please wait.'; return
    }
    sendStopTyping()
    chatState.socket.send(JSON.stringify({ type:'message', receiver_id:activeDmUserId.value, content }))
  } else if (activeType.value === 'group') {
    if (!groupSocket || groupSocket.readyState !== WebSocket.OPEN) {
      chatError.value = 'Connecting… please try again in a moment.'; return
    }
    groupSocket.send(JSON.stringify({ group_id:activeGroupId.value, content }))
    chatError.value = ''
  }
  newMsg.value = ''; typingTimeout = null; chatState.typingFrom = null
}

async function scrollBottom() { await nextTick(); if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight }

// ── Load older messages on scroll-to-top ─────────────────────────────────────
async function onChatScroll() {
  if (!chatBox.value || loadingOlder.value || noMoreMsgs.value || _scrollCooldown) return
  if (chatBox.value.scrollTop > 60) return // Only trigger near the very top

  const firstMsg = messages.value[0]
  if (!firstMsg) return

  const before = firstMsg.created_at
  loadingOlder.value = true
  _scrollCooldown = true
  const prevScrollHeight = chatBox.value.scrollHeight

  try {
    let older = []
    if (activeType.value === 'dm' && activeDmUserId.value) {
      const d = await apiFetch(`/messages/private?user_id=${activeDmUserId.value}&before=${encodeURIComponent(before)}`)
      older = Array.isArray(d) ? d : []
    } else if (activeType.value === 'group' && activeGroupId.value) {
      const d = await apiFetch(`/messages/group?group_id=${activeGroupId.value}&before=${encodeURIComponent(before)}`)
      older = Array.isArray(d) ? d : []
    }

    if (older.length === 0) {
      noMoreMsgs.value = true
    } else {
      messages.value = [...older, ...messages.value]
      await nextTick()
      if (chatBox.value) {
        chatBox.value.scrollTop = chatBox.value.scrollHeight - prevScrollHeight
      }
    }
  } catch {}
  finally {
    loadingOlder.value = false
    // 500 ms cooldown before the next page can be triggered
    setTimeout(() => { _scrollCooldown = false }, 500)
  }
}

// ── Global WS events (DM only) ───────────────────────────
let unsub = null
function handleWsEvent(msg) {
  if (msg.type === 'message') {
    // If it is for the currently open DM conversation, push to messages
    if (activeType.value === 'dm' &&
        ((msg.sender_id===activeDmUserId.value && msg.receiver_id===currentUserId.value) ||
         (msg.sender_id===currentUserId.value && msg.receiver_id===activeDmUserId.value))) {
      messages.value.push(msg); scrollBottom()
      if (msg.sender_id !== currentUserId.value) {
        clearUnread(msg.sender_id)
        apiFetch('/messages/mark-read', { method:'POST', body:JSON.stringify({ sender_id:msg.sender_id }) }).catch(()=>{})
        if (chatState.socket?.readyState===WebSocket.OPEN)
          chatState.socket.send(JSON.stringify({ type:'read', sender_id:msg.sender_id }))
      }
    }

    // ALWAYS update recent convo preview in sidebar and move it to the top!
    const peerId = msg.sender_id===currentUserId.value ? msg.receiver_id : msg.sender_id
    let u = users.value.find(u=>u.id===peerId)
    if (!u) {
      // Find details from profile endpoint
      apiFetch(`/profile?user_id=${peerId}`).then(prof => {
        if (prof) {
          const newUser = {
            id: peerId,
            first_name: prof.first_name,
            last_name: prof.last_name,
            avatar: prof.avatar || '',
            nickname: prof.nickname || ''
          }
          users.value.push(newUser)
          const idx = recentConvos.value.findIndex(c=>c.user.id===peerId)
          if (idx === -1) {
            recentConvos.value.unshift({ user: newUser, lastMessage: msg.content, lastAt: msg.created_at || new Date().toISOString() })
          }
        }
      }).catch(()=>{})
    } else {
      const idx = recentConvos.value.findIndex(c=>c.user.id===peerId)
      if (idx !== -1) {
        const c = recentConvos.value.splice(idx, 1)[0]
        c.lastMessage = msg.content
        c.lastAt = msg.created_at || new Date().toISOString()
        recentConvos.value.unshift(c)
      } else {
        recentConvos.value.unshift({ user:u, lastMessage:msg.content, lastAt:msg.created_at || new Date().toISOString() })
      }
    }
  }
  if (msg.type === 'read' && activeType.value === 'dm') {
    messages.value.forEach(m => { if (m.sender_id===currentUserId.value) m.is_read = true })
  }
}

async function loadOnline() {
  try { const d=await apiFetch('/users/online'); if(Array.isArray(d)) onlineIds.value=new Set(d) } catch {}
}

function startChatWith(user) {
  showNewChatModal.value = false
  modalSearch.value = ''
  if (!users.value.some(u => u.id === user.id)) {
    users.value.push(user)
  }
  openDm(user)
}

watch(() => route.query.user_id, async (newQueryUserId) => {
  if (newQueryUserId) {
    let qUser = users.value.find(u => u.id === newQueryUserId)
    if (!qUser) {
      try {
        const prof = await apiFetch(`/profile?user_id=${newQueryUserId}`)
        if (prof) {
          qUser = {
            id: newQueryUserId,
            first_name: prof.first_name,
            last_name: prof.last_name,
            avatar: prof.avatar || '',
            nickname: prof.nickname || ''
          }
          users.value.push(qUser)
        }
      } catch {}
    }
    if (qUser) {
      openDm(qUser)
    }
  }
})

onMounted(async () => {
  try { const me=await apiFetch('/me'); currentUserId.value=me.id; setSelfUserId(me.id) } catch { router.push('/login'); return }
  // Load both following and followers so either side can start a chat
  try {
    const [following, followers] = await Promise.all([
      apiFetch('/following?user_id='+currentUserId.value),
      apiFetch('/followers?user_id='+currentUserId.value),
    ])
    if (Array.isArray(following)) {
      followingUsers.value = following.filter(u => u.id !== currentUserId.value)
    }
    const map = new Map()
    if (Array.isArray(following)) {
      following.forEach(u => map.set(u.id, { id:u.id, first_name:u.first_name, last_name:u.last_name, avatar:u.avatar||'', nickname:u.nickname||'' }))
    }
    if (Array.isArray(followers)) {
      followers.forEach(u => {
        const id = u.user_id || u.id
        if (!map.has(id)) map.set(id, { id, first_name:u.first_name, last_name:u.last_name, avatar:u.avatar||'', nickname:u.nickname||'' })
      })
    }
    users.value = [...map.values()].filter(u => u.id !== currentUserId.value)
  } catch {}
  try { const d=await apiFetch('/groups/my'); myGroups.value=Array.isArray(d)?d.map(g=>({ id:g.id, title:g.title, member_count:g.member_count||0 })):[] } catch {}
  try {
    const recent = await apiFetch('/messages/recent')
    if (Array.isArray(recent)) {
      recentConvos.value = recent.map(c => ({
        user: {
          id: c.id,
          first_name: c.first_name,
          last_name: c.last_name,
          avatar: c.avatar,
          nickname: c.nickname
        },
        lastMessage: c.last_message,
        lastAt: c.last_message_time
      }))
    }
  } catch {}
  try {
    const unreads=await apiFetch('/messages/unread-count')
    if (Array.isArray(unreads)) {
      chatState.unreadTotal = 0
      chatState.unreadsBySender = {}
      for (const uc of unreads) {
        chatState.unreadsBySender[uc.sender_id] = uc.count
        chatState.unreadTotal += uc.count
      }
    }
  } catch {}
  unsub = onChatEvent(handleWsEvent)
  loadOnline()

  // Handle initial user_id from query parameter
  const queryUserId = route.query.user_id
  if (queryUserId) {
    let qUser = users.value.find(u => u.id === queryUserId)
    if (!qUser) {
      try {
        const prof = await apiFetch(`/profile?user_id=${queryUserId}`)
        if (prof) {
          qUser = {
            id: queryUserId,
            first_name: prof.first_name,
            last_name: prof.last_name,
            avatar: prof.avatar || '',
            nickname: prof.nickname || ''
          }
          users.value.push(qUser)
        }
      } catch {}
    }
    if (qUser) {
      openDm(qUser)
    }
  }
})

onUnmounted(() => {
  if (unsub) unsub()
  closeGroupSocket()
  globalActiveChatUserId.value = null
  globalActiveChatGroupId.value = null
  // Clear typing indicator so it doesn't persist after leaving chat
  chatState.typingFrom = null
  clearTimeout(chatState._typingTimer)
})
</script>

<style scoped>
.sidebar-item {
  display:flex; align-items:center; gap:10px;
  padding:8px 10px; border-radius:10px;
  cursor:pointer; transition:background 0.14s; margin-bottom:2px;
}
.sidebar-item:hover, .sidebar-item.active { background:var(--surface2); }
.sidebar-item.active { background:var(--accent-bg); }

.presence-dot {
  position:absolute; bottom:0; right:0;
  width:9px; height:9px; border-radius:50%;
  background:var(--success); border:2px solid var(--surface);
}

.status-online, .status-offline { width:7px; height:7px; border-radius:50%; display:inline-block; }
.status-online { background:var(--success); }
.status-offline { background:var(--text-3); }

.typing-dots { display:inline-flex; align-items:center; gap:3px; }
.typing-dots span { width:5px; height:5px; background:var(--text-3); border-radius:50%; animation:tdot 1.2s infinite; }
.typing-dots span:nth-child(2) { animation-delay:0.2s; }
.typing-dots span:nth-child(3) { animation-delay:0.4s; }
@keyframes tdot {
  0%,60%,100% { transform:translateY(0);opacity:0.6; }
  30% { transform:translateY(-4px);opacity:1; }
}
</style>
