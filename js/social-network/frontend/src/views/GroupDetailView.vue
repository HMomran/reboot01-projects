<template>
  <div>
    <Navbar />
    
    <div v-if="groupInfo" class="page-layout" style="max-width:1100px; margin: 24px auto; display: flex; flex-direction: column; gap: 20px; padding: 0 16px">
      <!-- Group info card (visible to everyone) -->
      <div class="card" style="padding:22px">
        <div style="display:flex;align-items:flex-start;justify-content:space-between;gap:10px;margin-bottom:12px">
          <div>
            <h2 style="font-size:1.20rem;margin-bottom:4px;font-weight:600">{{ groupInfo?.title || 'Group' }}</h2>
            <p style="font-size:0.83rem;color:var(--text-2)">{{ groupInfo?.description }}</p>
          </div>
          <div style="display:flex;gap:6px" v-if="isMember">
            <button class="btn btn-ghost btn-sm" style="display:inline-flex;align-items:center;gap:4px" @click="showMembers=true">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
              Members
            </button>
            <button class="btn btn-danger btn-sm" @click="leaveGroup">Leave</button>
          </div>
          <div v-else style="display:flex;gap:6px;align-items:center">
            <button v-if="groupInfo?.membership_status === ''" class="btn btn-primary btn-sm" @click="joinGroup">Request to Join</button>
            <span v-else-if="groupInfo?.membership_status === 'requested'" style="font-size:0.85rem;color:var(--text-3);font-style:italic">Request Pending</span>
            <template v-else-if="groupInfo?.membership_status === 'invited'">
              <span style="font-size:0.85rem;color:var(--accent);font-weight:600">You're invited!</span>
              <button class="btn btn-primary btn-sm" @click="acceptInviteFromDetail">Accept</button>
              <button class="btn btn-ghost btn-sm" @click="declineInviteFromDetail">Decline</button>
            </template>
          </div>
        </div>

        <!-- Add Member -->
        <div style="display:flex;gap:8px;margin-top:14px" v-if="isMember">
          <select class="input" v-model="inviteUserId" style="flex:1;padding:8px 10px;font-size:0.83rem">
            <option value="">Invite member…</option>
            <option v-for="u in nonMembers" :key="u.id" :value="u.id">{{ u.first_name }} {{ u.last_name }}</option>
          </select>
          <button class="btn btn-primary btn-sm" :disabled="!inviteUserId" @click="addMember">Invite</button>
        </div>
      </div>

      <!-- If NOT a member -->
      <div v-if="!isMember" class="card" style="padding: 48px; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 16px">
        <div style="width: 60px; height: 60px; border-radius: 50%; background: var(--surface2); display: flex; align-items: center; justify-content: center; color: var(--text-3)">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
        </div>
        <h3 style="font-size: 1.1rem; font-weight: 600">This group is private</h3>
        <p style="font-size: 0.88rem; color: var(--text-2); max-width: 400px; margin: 0">You must be an accepted member of this group to view its posts, events, and member list.</p>
        <button v-if="groupInfo.membership_status === ''" class="btn btn-primary" style="margin-top: 8px" @click="joinGroup">Request to Join</button>
        <span v-else-if="groupInfo.membership_status === 'requested'" style="font-size: 0.9rem; color: var(--text-3); font-style: italic; margin-top: 8px">Request Pending Approval</span>
        <div v-else-if="groupInfo.membership_status === 'invited'" style="display:flex;gap:10px;margin-top:8px">
          <button class="btn btn-primary" @click="acceptInviteFromDetail">Accept Invitation</button>
          <button class="btn btn-ghost" @click="declineInviteFromDetail">Decline</button>
        </div>
      </div>

      <!-- If IS a member, render the regular layout -->
      <div v-else class="two-col" style="display: grid; grid-template-columns: 1fr 2fr; gap: 20px">
        <!-- Left column: Events -->
        <div style="display:flex;flex-direction:column;gap:14px">
          <!-- Events -->
          <div class="card" style="padding:22px">
            <h3 style="font-size:1rem;margin-bottom:14px;display:flex;align-items:center;gap:6px">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
              Events
            </h3>

            <!-- Create event -->
            <div style="background:var(--surface2);border-radius:10px;padding:14px;margin-bottom:14px">
              <div class="form-field">
                <label>Title</label>
                <input class="input" v-model="newEvent.title" placeholder="Event name" style="font-size:0.85rem" :style="{ borderColor: newEvent.title.length > 100 ? 'var(--danger)' : '' }" />
                <div style="text-align:right;font-size:0.68rem;margin-top:2px" :style="{ color: newEvent.title.length > 100 ? 'var(--danger)' : 'var(--text-3)' }">{{ newEvent.title.length }}/100</div>
              </div>
              <div class="form-field">
                <label>Description</label>
                <textarea class="input" v-model="newEvent.description" rows="2" style="font-size:0.85rem" :style="{ borderColor: newEvent.description.length > 500 ? 'var(--danger)' : '' }"></textarea>
                <div style="text-align:right;font-size:0.68rem;margin-top:2px" :style="{ color: newEvent.description.length > 500 ? 'var(--danger)' : 'var(--text-3)' }">{{ newEvent.description.length }}/500</div>
              </div>
              <div class="form-field"><label>Date & Time</label><input class="input" v-model="newEvent.event_datetime" type="datetime-local" style="font-size:0.85rem" /></div>
              <div v-if="eventCreateErr" class="error-msg" style="margin-bottom:8px;padding:6px 10px">{{ eventCreateErr }}</div>
              <button class="btn btn-primary btn-sm w-full" :disabled="!newEvent.title.trim()||!newEvent.description.trim()||!newEvent.event_datetime||newEvent.title.length>100||newEvent.description.length>500||creatingEvent" @click="createEvent">
                {{ creatingEvent ? '…' : 'Create Event' }}
              </button>
            </div>

            <div v-if="!events.length" style="color:var(--text-3);font-size:0.83rem">No events yet.</div>

            <div v-for="ev in events" :key="ev.id" style="border:1px solid var(--border);border-radius:10px;padding:14px;margin-bottom:8px">
              <div style="font-weight:700;font-size:0.9rem;margin-bottom:2px">{{ ev.title }}</div>
              <div style="font-size:0.78rem;color:var(--text-2);margin-bottom:6px">{{ ev.description }}</div>
              <div style="font-size:0.75rem;color:var(--accent);margin-bottom:10px;display:inline-flex;align-items:center;gap:4px">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                {{ formatDate(ev.event_datetime) }}
              </div>

              <!-- Responses summary -->
              <div v-if="eventResponses[ev.id]" style="font-size:0.76rem;color:var(--text-3);margin-bottom:8px;display:flex;align-items:center;gap:8px">
                <span style="display:inline-flex;align-items:center;gap:3px;color:var(--success)"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> {{ goingCount(ev.id) }} going</span>
                <span>·</span>
                <span style="display:inline-flex;align-items:center;gap:3px;color:var(--danger)"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg> {{ notGoingCount(ev.id) }} not going</span>
              </div>

              <div style="display:flex;gap:6px;flex-wrap:wrap">
                <button class="btn btn-primary btn-sm" @click="respond(ev.id,'going')" style="display:inline-flex;align-items:center;gap:4px">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  Going
                </button>
                <button class="btn btn-ghost btn-sm" @click="respond(ev.id,'not_going')" style="display:inline-flex;align-items:center;gap:4px">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  Not going
                </button>
                <button class="btn btn-ghost btn-sm" @click="loadResponses(ev.id)">Show responses</button>
              </div>

              <div v-if="eventResponses[ev.id]?.length" style="margin-top:10px;display:flex;flex-direction:column;gap:4px">
                <div v-for="r in eventResponses[ev.id]" :key="r.user_id+r.created_at" style="font-size:0.75rem;color:var(--text-2)">
                  {{ r.user_id === myId ? 'You' : (r.first_name + ' ' + r.last_name) }} —
                  <span style="display:inline-flex;align-items:center;gap:2px" :style="{ color: r.response === 'going' ? 'var(--success)' : 'var(--danger)' }">
                    <svg v-if="r.response === 'going'" xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                    {{ r.response }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Join requests (creator only) -->
          <div class="card" style="padding:22px" v-if="joinRequests.length">
            <h3 style="font-size:1rem;margin-bottom:14px;display:flex;align-items:center;gap:6px">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
              Join Requests
            </h3>
            <div v-for="req in joinRequests" :key="req.user_id" style="display:flex;align-items:center;gap:10px;margin-bottom:10px">
              <div class="avatar-ring" style="width:32px;height:32px;font-size:0.68rem;flex-shrink:0" :style="{ background: avatarColor(req.user_id) }">{{ initials(req.first_name,req.last_name) }}</div>
              <span style="flex:1;font-size:0.88rem">{{ req.first_name }} {{ req.last_name }}</span>
              <button class="btn btn-secondary btn-sm" style="display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;color:var(--accent);border-color:var(--accent);background:rgba(0,122,253,0.08);padding:0" @click="acceptReq(req.user_id)" title="Accept">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              </button>
              <button class="btn btn-secondary btn-sm" style="display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;color:var(--danger);border-color:var(--danger);background:rgba(224,64,64,0.08);padding:0" @click="declineReq(req.user_id)" title="Decline">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Right column: Posts -->
        <div>
          <!-- Compose post -->
          <div class="card" style="padding:18px;margin-bottom:14px">
            <textarea class="input" v-model="newPostContent" rows="3" placeholder="Write something for the group…" style="margin-bottom:4px;font-size:0.92rem"
              :style="{ borderColor: newPostContent.length > 2000 ? 'var(--danger)' : '' }"></textarea>
            <div style="text-align:right;font-size:0.7rem;margin-bottom:6px" :style="{ color: newPostContent.length > 2000 ? 'var(--danger)' : 'var(--text-3)' }">{{ newPostContent.length }}/2000</div>
            <div v-if="postErr" class="error-msg" style="margin-bottom:8px;padding:8px 12px">{{ postErr }}</div>
            <button class="btn btn-primary btn-sm" :disabled="!newPostContent.trim()||newPostContent.length>2000||postingPost" @click="createPost">
              {{ postingPost ? 'Posting…' : 'Post to Group' }}
            </button>
          </div>

          <!-- Posts -->
          <div v-if="!posts.length" class="card empty-state" style="padding:48px">
            <div class="empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            </div>
            <h3>No posts yet</h3>
            <p>Be the first to post!</p>
          </div>

          <article class="post-card" style="margin-bottom:12px" v-for="post in posts" :key="post.id">
            <div class="post-header">
              <div class="avatar-ring post-avatar" :style="{ background: avatarColor(post.user_id) }">
                <img v-if="post.avatar" :src="BACKEND_URL + post.avatar" />
                <span v-else>{{ initials(post.first_name, post.last_name) }}</span>
              </div>
              <div class="post-meta">
                <span class="post-author">{{ post.user_id===myId ? 'You' : post.first_name+' '+post.last_name }}</span>
                <span class="post-time">{{ relTime(post.created_at) }}</span>
              </div>
              <button v-if="post.user_id===myId || groupInfo?.creator_id===myId" class="btn btn-danger btn-sm" style="margin-left:auto;display:inline-flex;align-items:center;justify-content:center;padding:6px 8px" @click="deleteGroupPost(post)">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
              </button>
            </div>
            <p class="post-body" v-if="post.content">{{ post.content }}</p>
            <img v-if="post.image" :src="BACKEND_URL + post.image" class="post-img" />

            <div class="post-actions">
              <button class="react-btn" @click="toggleGroupComments(post)" style="display:inline-flex;align-items:center;gap:6px">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
                {{ (groupCommentStates[post.id]?.comments||[]).length }}
              </button>
            </div>

            <!-- Comments -->
            <div class="comments-section" v-if="groupCommentStates[post.id]?.open">
              <div v-for="c in groupCommentStates[post.id].comments" :key="c.id" class="comment-item">
                <div class="avatar-ring" style="width:28px;height:28px;font-size:0.65rem;flex-shrink:0" :style="{ background: avatarColor(c.user_id) }">
                  <img v-if="c.avatar" :src="BACKEND_URL + c.avatar" />
                  <span v-else>{{ initials(c.first_name, c.last_name) }}</span>
                </div>
                <div class="comment-bubble" style="flex:1">
                  <span class="comment-author">{{ c.first_name }} {{ c.last_name }}</span>
                  <p class="comment-text">{{ c.content }}</p>
                  <span class="comment-time">{{ relTime(c.created_at) }}</span>
                </div>
                <button v-if="c.user_id===myId" class="btn btn-danger btn-sm" style="flex-shrink:0;align-self:flex-start;display:inline-flex;align-items:center;justify-content:center;padding:6px" @click="deleteGroupComment(post, c)">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
                </button>
              </div>
              <div class="comment-composer">
                <input class="input" v-model="groupCommentStates[post.id].text" placeholder="Comment…" @keyup.enter="addGroupComment(post)"
                  :style="{ borderColor: (groupCommentStates[post.id].text?.length||0) > 500 ? 'var(--danger)' : '' }" />
                <span style="font-size:0.68rem;white-space:nowrap" :style="{ color: (groupCommentStates[post.id].text?.length||0) > 500 ? 'var(--danger)' : 'var(--text-3)' }">
                  {{ groupCommentStates[post.id].text?.length||0 }}/500
                </span>
                <button class="btn btn-primary btn-sm" :disabled="!groupCommentStates[post.id].text?.trim()||(groupCommentStates[post.id].text?.length||0)>500" @click="addGroupComment(post)">Send</button>
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>

    <!-- Group not found fallback -->
    <div v-else class="page-layout" style="max-width:600px; margin: 40px auto; text-align: center">
      <div class="card" style="padding: 40px">
        <h3 style="font-weight:600; margin-bottom:12px">Group not found</h3>
        <p style="color: var(--text-3); font-size: 0.9rem">This group does not exist or has been dissolved.</p>
        <router-link to="/groups" class="btn btn-primary" style="margin-top: 20px; display: inline-block">Back to Groups</router-link>
      </div>
    </div>

    <!-- Members modal -->
    <div v-if="showMembers" class="modal-backdrop" @click.self="showMembers=false">
      <div class="modal-box">
        <div class="modal-header"><h3>Members</h3><button class="modal-close" @click="showMembers=false">✕</button></div>
        <div v-for="m in members" :key="m.id" style="display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid var(--border)">
          <div class="avatar-ring" style="width:34px;height:34px;font-size:0.72rem;flex-shrink:0" :style="{ background: avatarColor(m.id) }">
            <img v-if="m.avatar" :src="BACKEND_URL + m.avatar" />
            <span v-else>{{ initials(m.first_name, m.last_name) }}</span>
          </div>
          <span style="font-size:0.9rem;font-weight:500;flex:1">{{ m.id===myId ? 'You' : m.first_name+' '+m.last_name }}</span>
          <span v-if="m.is_creator" style="font-size:0.72rem;background:var(--accent-bg);color:var(--accent);padding:2px 8px;border-radius:20px;font-weight:600">Admin</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'

const route = useRoute(); const router = useRouter()
const groupId = route.params.id

const groupInfo = ref(null); const posts = ref([]); const events = ref([])
const members = ref([]); const joinRequests = ref([]); const allUsers = ref([])
const eventResponses = reactive({}); const groupCommentStates = reactive({})
const myId = ref(''); const showMembers = ref(false); const inviteUserId = ref('')
const newPostContent = ref(''); const postingPost = ref(false); const postErr = ref('')
const newEvent = reactive({ title:'', description:'', event_datetime:'' }); const creatingEvent = ref(false); const eventCreateErr = ref('')
const error = ref('')

const isMember = computed(() => {
  return groupInfo.value?.membership_status === 'accepted'
})

const isCreator = computed(() => {
  return groupInfo.value?.creator_id === myId.value
})

const nonMembers = computed(() => allUsers.value.filter(u => !members.value.some(m => m.id === u.id) && u.id !== myId.value))

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
function formatDate(v) { if (!v) return ''; return v.replace('T',' ').substring(0,16) }
function goingCount(id) { return (eventResponses[id]||[]).filter(r=>r.response==='going').length }
function notGoingCount(id) { return (eventResponses[id]||[]).filter(r=>r.response==='not_going').length }

function ensureGCS(id) { if (!groupCommentStates[id]) groupCommentStates[id]={open:false,loading:false,comments:[],text:''} }

async function loadPosts() {
  try { const d=await apiFetch(`/groups/posts?group_id=${groupId}`); posts.value=Array.isArray(d)?d:[]; posts.value.forEach(p=>ensureGCS(p.id)) } catch {}
}
async function loadEvents() {
  try { const d=await apiFetch(`/groups/events?group_id=${groupId}`); events.value=Array.isArray(d)?d:[] } catch {}
}
async function loadMembers() {
  try { const d=await apiFetch(`/groups/members?group_id=${groupId}`); members.value=Array.isArray(d)?d:[] } catch {}
}
async function loadJoinRequests() {
  try { const d=await apiFetch(`/groups/requests?group_id=${groupId}`); joinRequests.value=Array.isArray(d)?d:[] } catch {}
}

async function loadGroupInfo() {
  try {
    const d=await apiFetch('/groups')
    if (Array.isArray(d)) groupInfo.value=d.find(g=>g.id===groupId)||null
  } catch {}
}

async function createPost() {
  postingPost.value=true; postErr.value=''
  try {
    await apiFetch('/groups/posts/create', { method:'POST', body:JSON.stringify({ group_id:groupId, content:newPostContent.value }) })
    newPostContent.value=''; await loadPosts()
  } catch(e) { postErr.value=e.message } finally { postingPost.value=false }
}

async function createEvent() {
  eventCreateErr.value = ''
  if (!newEvent.title.trim() || !newEvent.description.trim() || !newEvent.event_datetime) {
    eventCreateErr.value = 'All event fields are required.'; return
  }
  if (newEvent.title.trim().length > 100) { eventCreateErr.value = 'Title cannot exceed 100 characters.'; return }
  if (newEvent.description.trim().length > 500) { eventCreateErr.value = 'Description cannot exceed 500 characters.'; return }
  const dt = new Date(newEvent.event_datetime)
  if (isNaN(dt.getTime()) || dt <= new Date()) { eventCreateErr.value = 'Event date must be in the future.'; return }
  creatingEvent.value=true
  try {
    await apiFetch('/groups/events/create', { method:'POST', body:JSON.stringify({ group_id:groupId, ...newEvent }) })
    newEvent.title=''; newEvent.description=''; newEvent.event_datetime=''; eventCreateErr.value=''; await loadEvents()
  } catch(e) { eventCreateErr.value=e.message } finally { creatingEvent.value=false }
}

async function respond(eventId, response) {
  try { await apiFetch('/groups/events/respond', { method:'POST', body:JSON.stringify({ event_id:eventId, response }) }); await loadResponses(eventId) } catch {}
}

async function loadResponses(eventId) {
  try { const d=await apiFetch(`/groups/events/responses?event_id=${eventId}`); eventResponses[eventId]=Array.isArray(d)?d:[] } catch {}
}

async function addMember() {
  if (!inviteUserId.value) return
  const invitedUser = allUsers.value.find(u => u.id === inviteUserId.value)
  const invitedName = invitedUser ? `${invitedUser.first_name} ${invitedUser.last_name}` : 'User'
  try {
    await apiFetch('/groups/invite', { method:'POST', body:JSON.stringify({ group_id:groupId, user_id:inviteUserId.value }) })
    inviteUserId.value=''
    alert(`Invitation sent to ${invitedName}. They will appear as a member once they accept.`)
    await loadMembers()
  } catch(e) { error.value=e.message }
}

async function leaveGroup() {
  if (!confirm('Leave this group?')) return
  try { await apiFetch(`/groups/leave?group_id=${groupId}`, { method:'POST' }); router.push('/groups') } catch(e) { error.value=e.message }
}

async function joinGroup() {
  try {
    await apiFetch('/groups/request', { method: 'POST', body: JSON.stringify({ group_id: groupId }) })
    if (groupInfo.value) {
      groupInfo.value.membership_status = 'requested'
    }
  } catch(e) { error.value=e.message }
}

async function acceptInviteFromDetail() {
  try {
    await apiFetch('/groups/invite/accept', { method:'POST', body:JSON.stringify({ group_id:groupId }) })
    await loadGroupInfo()
    if (isMember.value) {
      await Promise.all([loadPosts(), loadEvents(), loadMembers(), loadJoinRequests()])
    }
  } catch(e) { error.value=e.message }
}

async function declineInviteFromDetail() {
  try {
    await apiFetch('/groups/invite/decline', { method:'POST', body:JSON.stringify({ group_id:groupId }) })
    router.push('/groups')
  } catch(e) { error.value=e.message }
}

async function acceptReq(userId) {
  try { await apiFetch('/groups/request/accept', { method:'POST', body:JSON.stringify({ group_id:groupId, user_id:userId }) }); await loadJoinRequests(); await loadMembers() } catch {}
}
async function declineReq(userId) {
  try { await apiFetch('/groups/request/decline', { method:'POST', body:JSON.stringify({ group_id:groupId, user_id:userId }) }); await loadJoinRequests() } catch {}
}

async function deleteGroupPost(post) {
  if (!confirm('Delete this post?')) return
  try {
    await apiFetch(`/groups/posts/delete?post_id=${post.id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== post.id)
  } catch (e) { alert(e.message || 'Failed to delete post.') }
}

async function deleteGroupComment(post, comment) {
  if (!confirm('Delete this comment?')) return
  const st = groupCommentStates[post.id]
  try {
    await apiFetch(`/groups/posts/comments/delete?comment_id=${comment.id}`, { method: 'DELETE' })
    st.comments = st.comments.filter(c => c.id !== comment.id)
  } catch (e) { alert(e.message || 'Failed to delete comment.') }
}

async function toggleGroupComments(post) {
  ensureGCS(post.id); const st=groupCommentStates[post.id]
  st.open=!st.open
  if (st.open && !st.comments.length) {
    st.loading=true
    try { const d=await apiFetch(`/groups/posts/comments?post_id=${post.id}`); st.comments=Array.isArray(d)?d:[] }
    catch {} finally { st.loading=false }
  }
}

async function addGroupComment(post) {
  const st=groupCommentStates[post.id]; if (!st.text?.trim()) return
  try {
    await apiFetch('/groups/posts/comments/create', { method:'POST', body:JSON.stringify({ post_id:post.id, content:st.text.trim() }) })
    st.text=''
    const d=await apiFetch(`/groups/posts/comments?post_id=${post.id}`); st.comments=Array.isArray(d)?d:[]
  } catch {}
}

onMounted(async () => {
  try { const me=await apiFetch('/me'); myId.value=me.id } catch { router.push('/login'); return }
  try { const d=await apiFetch('/users'); allUsers.value=Array.isArray(d)?d:[] } catch {}
  await loadGroupInfo()
  if (isMember.value) {
    await Promise.all([loadPosts(), loadEvents(), loadMembers(), loadJoinRequests()])
  }
})
</script>
