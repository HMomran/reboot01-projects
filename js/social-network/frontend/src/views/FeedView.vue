<template>
  <div>
    <Navbar />
    <div class="page-layout feed-col">

      <!-- Composer -->
      <div v-if="isLoggedIn" class="card" style="padding:20px">
        <div style="display:flex;gap:12px;align-items:flex-start">
          <div class="avatar-ring" style="width:40px;height:40px;font-size:0.85rem;flex-shrink:0" :style="{ background: avatarColor(myId) }">
            <img v-if="myAvatar" :src="BACKEND_URL + myAvatar" />
            <span v-else>{{ myInitials }}</span>
          </div>
          <div style="flex:1">
            <textarea class="input" v-model="newContent" rows="2"
              placeholder="What's on your mind?"
              style="resize:none;border-radius:12px;font-size:0.92rem"
              @input="e => { e.target.style.height='auto'; e.target.style.height=e.target.scrollHeight+'px' }"
              :style="{ borderColor: newContent.length > 2000 ? 'var(--danger)' : '' }"
            ></textarea>
            <div style="text-align:right;font-size:0.72rem;margin-top:3px" :style="{ color: newContent.length > 2000 ? 'var(--danger)' : 'var(--text-3)' }">{{ newContent.length }}/2000</div>

            <!-- Image preview -->
            <div v-if="imagePreview" style="position:relative;display:inline-block;margin-top:10px">
              <img :src="imagePreview" style="max-height:160px;border-radius:10px;border:1px solid var(--border)" />
              <button @click="removeImage" style="position:absolute;top:6px;right:6px;background:rgba(0,0,0,0.65);color:#fff;border:none;border-radius:50%;width:22px;height:22px;cursor:pointer;display:flex;align-items:center;justify-content:center;padding:0">
                <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>

            <!-- Privacy + allowed users -->
            <div style="margin-top:12px;display:flex;flex-direction:column;gap:10px">
              <div style="display:flex;align-items:center;justify-content:space-between;gap:8px;flex-wrap:wrap">
                <div style="display:flex;align-items:center;gap:8px">
                  <label style="cursor:pointer;color:var(--text-2);font-size:0.85rem;display:flex;align-items:center;gap:5px;padding:6px 10px;background:var(--surface2);border-radius:8px;border:1px solid var(--border)">
                    <input type="file" accept="image/*" style="display:none" ref="fileRef" @change="pickImage" />
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"></path></svg>
                    Photo
                  </label>
                  <select class="input" v-model="privacy" style="width:auto;padding:6px 10px;font-size:0.82rem">
                    <option value="public">Public</option>
                    <option value="followers">Followers only</option>
                    <option value="private">Selected followers</option>
                  </select>
                </div>
                <button class="btn btn-primary btn-sm" :disabled="(!newContent.trim() && !selectedImage) || newContent.length > 2000 || posting" @click="createPost">
                  {{ posting ? 'Posting…' : 'Post' }}
                </button>
              </div>

              <!-- Selected followers picker (only when privacy=private) -->
              <div v-if="privacy === 'private'" style="background:var(--surface2);border-radius:10px;padding:12px;border:1px solid var(--border)">
                <div style="font-size:0.78rem;font-weight:600;color:var(--text-2);text-transform:uppercase;letter-spacing:0.06em;margin-bottom:8px">Who can see this post?</div>
                <div v-if="!myFollowers.length" style="font-size:0.82rem;color:var(--text-3)">You have no followers yet.</div>
                <div style="display:flex;flex-wrap:wrap;gap:6px">
                  <label v-for="f in myFollowers" :key="f.user_id"
                    :style="{ display:'flex', alignItems:'center', gap:'6px', padding:'5px 10px', borderRadius:'20px', cursor:'pointer', fontSize:'0.82rem', fontWeight:'500', border:'1px solid', borderColor: allowedUsers.includes(f.user_id) ? 'var(--accent)' : 'var(--border)', background: allowedUsers.includes(f.user_id) ? 'var(--accent-bg)' : 'var(--surface3)', color: allowedUsers.includes(f.user_id) ? 'var(--accent)' : 'var(--text-2)' }">
                    <input type="checkbox" :value="f.user_id" v-model="allowedUsers" style="display:none" />
                    <span style="display:inline-flex;align-items:center;gap:4px">
                      <svg v-if="allowedUsers.includes(f.user_id)" xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
                    </span>
                    {{ f.first_name }} {{ f.last_name }}
                  </label>
                </div>
              </div>
            </div>

            <div v-if="postError" class="error-msg" style="margin-top:8px;padding:8px 12px">{{ postError }}</div>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <template v-if="loading">
        <div class="card" style="padding:20px" v-for="n in 3" :key="n">
          <div style="display:flex;gap:12px;margin-bottom:14px">
            <div class="skeleton" style="width:40px;height:40px;border-radius:50%;flex-shrink:0"></div>
            <div style="flex:1;display:flex;flex-direction:column;gap:6px;padding-top:4px">
              <div class="skeleton" style="height:12px;width:38%"></div>
              <div class="skeleton" style="height:10px;width:22%"></div>
            </div>
          </div>
          <div class="skeleton" style="height:13px;width:100%;margin-bottom:8px"></div>
          <div class="skeleton" style="height:13px;width:70%"></div>
        </div>
      </template>

      <!-- Error -->
      <div v-else-if="fetchError" class="card empty-state">
        <p style="display:inline-flex;align-items:center;gap:6px;justify-content:center;color:var(--danger)">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
          {{ fetchError }}
        </p>
        <button class="btn btn-secondary btn-sm" style="margin-top:12px" @click="loadFeed">Retry</button>
      </div>

      <!-- Empty -->
      <div v-else-if="posts.length === 0" class="card empty-state">
        <div class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
        </div>
        <h3>Your feed is empty</h3>
        <p>Follow people or create your first post!</p>
        <router-link to="/explore" class="btn btn-primary btn-sm" style="margin-top:16px;display:inline-flex">Explore →</router-link>
      </div>

      <!-- Posts -->
      <article class="post-card" v-for="post in posts" :key="post.id">
        <div class="post-header">
          <router-link :to="'/profile/' + post.user_id">
            <div class="avatar-ring post-avatar" :style="{ background: avatarColor(post.user_id) }">
              <img v-if="post.avatar" :src="BACKEND_URL + post.avatar" />
              <span v-else>{{ initials(post.first_name, post.last_name) }}</span>
            </div>
          </router-link>
          <div class="post-meta">
            <router-link :to="'/profile/' + post.user_id" class="post-author">
              {{ post.user_id === myId ? 'You' : post.first_name + ' ' + post.last_name }}
            </router-link>
            <span class="post-time">{{ relTime(post.created_at) }}</span>
          </div>
          <span class="privacy-tag" style="display:inline-flex;align-items:center;gap:4px">
            <svg v-if="post.privacy === 'public'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
            <svg v-else-if="post.privacy === 'followers'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
            {{ post.privacy === 'public' ? 'Public' : post.privacy === 'followers' ? 'Followers' : 'Private' }}
          </span>
          <button v-if="post.user_id === myId" @click="deletePost(post)" title="Delete post"
            style="margin-left:auto;background:none;border:none;cursor:pointer;color:var(--text-3);font-size:0.85rem;padding:2px 6px;border-radius:6px;line-height:1;transition:color 0.15s;display:inline-flex;align-items:center;justify-content:center"
            @mouseenter="e=>e.currentTarget.style.color='var(--danger)'" @mouseleave="e=>e.currentTarget.style.color='var(--text-3)'">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
          </button>
        </div>

        <p class="post-body" v-if="post.content">{{ post.content }}</p>
        <img v-if="post.image" :src="BACKEND_URL + post.image" class="post-img" />

        <div class="post-actions" style="display:flex;gap:8px">
          <button :class="['react-btn', post.my_reaction === 'like' && 'active-like']" @click="isLoggedIn ? react(post,'like') : null" style="display:inline-flex;align-items:center;gap:6px" :style="!isLoggedIn ? { cursor: 'default', opacity: 0.75 } : {}">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
            {{ post.likes }}
          </button>
          <button :class="['react-btn', post.my_reaction === 'dislike' && 'active-dislike']" @click="isLoggedIn ? react(post,'dislike') : null" style="display:inline-flex;align-items:center;gap:6px" :style="!isLoggedIn ? { cursor: 'default', opacity: 0.75 } : {}">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zm7-13h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3"></path></svg>
            {{ post.dislikes }}
          </button>
          <button :class="['react-btn', states[post.id]?.open && 'active-comment']" @click="toggleComments(post)" style="display:inline-flex;align-items:center;gap:6px">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
            {{ states[post.id]?.open ? (states[post.id]?.comments || []).length : (states[post.id]?._pendingCommentCount ?? 0) }}
          </button>
        </div>

        <!-- Comments -->
        <div class="comments-section" v-if="states[post.id]?.open">
          <div v-if="states[post.id].loading" style="padding:8px 0">
            <div class="skeleton" style="height:36px;border-radius:10px;margin-bottom:8px"></div>
            <div class="skeleton" style="height:36px;border-radius:10px;width:80%"></div>
          </div>
          <div v-else-if="states[post.id].loadError" style="color:var(--danger);font-size:0.82rem;padding:8px 0">
            Failed to load comments. <button class="btn btn-ghost btn-sm" @click="toggleComments(post); toggleComments(post)">Retry</button>
          </div>
          <template v-else>
            <div v-if="!states[post.id].comments.length" style="color:var(--text-3);font-size:0.82rem;padding:4px 0 10px">No comments yet. Be first!</div>
            <div class="comment-item" v-for="c in states[post.id].comments" :key="c.id">
              <div class="avatar-ring" style="width:28px;height:28px;font-size:0.65rem;flex-shrink:0" :style="{ background: avatarColor(c.user_id) }">
                <img v-if="c.avatar" :src="BACKEND_URL + c.avatar" />
                <span v-else>{{ initials(c.first_name, c.last_name) }}</span>
              </div>
              <div class="comment-bubble">
                <span class="comment-author">{{ c.first_name }} {{ c.last_name }}</span>
                <p class="comment-text">{{ c.content }}</p>
                <img v-if="c.image" :src="BACKEND_URL + c.image" style="max-width:200px;border-radius:8px;margin-top:6px;display:block" />
                <div class="comment-footer">
                  <span class="comment-time">{{ relTime(c.created_at) }}</span>
                  <button :class="['comment-like-btn', c.i_liked && 'liked']" @click="isLoggedIn ? likeComment(post, c) : null" style="display:inline-flex;align-items:center;gap:4px" :style="!isLoggedIn ? { cursor: 'default', opacity: 0.75 } : {}">
                    <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                    {{ c.likes }}
                  </button>
                  <button v-if="c.user_id === myId || post.user_id === myId" @click="deleteComment(post, c)" title="Delete"
                    style="background:none;border:none;cursor:pointer;color:var(--text-3);font-size:0.75rem;padding:0 4px;line-height:1;transition:color 0.15s;display:inline-flex;align-items:center;justify-content:center"
                    @mouseenter="e=>e.currentTarget.style.color='var(--danger)'" @mouseleave="e=>e.currentTarget.style.color='var(--text-3)'">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
                  </button>
                </div>
              </div>
            </div>
          </template>

          <!-- Comment composer with image -->
          <div v-if="isLoggedIn" style="margin-top:10px">
            <!-- Comment image preview -->
            <div v-if="states[post.id].commentImagePreview" style="position:relative;display:inline-block;margin-bottom:8px">
              <img :src="states[post.id].commentImagePreview" style="max-height:100px;border-radius:8px;border:1px solid var(--border)" />
              <button @click="removeCommentImage(post.id)" style="position:absolute;top:4px;right:4px;background:rgba(0,0,0,0.65);color:#fff;border:none;border-radius:50%;width:20px;height:20px;cursor:pointer;display:flex;align-items:center;justify-content:center;padding:0">
                <svg xmlns="http://www.w3.org/2000/svg" width="9" height="9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>
            <div class="comment-composer">
              <label style="cursor:pointer;color:var(--text-3);padding:8px;display:flex;align-items:center">
                <input type="file" accept="image/*" style="display:none" @change="e => pickCommentImage(post.id, e)" />
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"></path></svg>
              </label>
              <input class="input" v-model="states[post.id].text" placeholder="Add a comment…" @keyup.enter="addComment(post)" style="flex:1"
                :style="{ borderColor: (states[post.id].text?.length||0) > 500 ? 'var(--danger)' : '' }" />
              <span style="font-size:0.7rem;white-space:nowrap" :style="{ color: (states[post.id].text?.length||0) > 500 ? 'var(--danger)' : 'var(--text-3)' }">{{ states[post.id].text?.length||0 }}/500</span>
              <button class="btn btn-primary btn-sm" :disabled="(!states[post.id].text?.trim() && !states[post.id].commentImage) || (states[post.id].text?.length||0) > 500 || states[post.id].sending" @click="addComment(post)">Send</button>
            </div>
            <div v-if="states[post.id].commentError" style="color:var(--danger);font-size:0.8rem;margin-top:6px">{{ states[post.id].commentError }}</div>
          </div>
        </div>
      </article>

    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'
import { isLoggedIn } from '../stores/auth'

const router = useRouter()
const posts = ref([])
const loading = ref(true)
const fetchError = ref('')
const posting = ref(false)
const postError = ref('')
const myId = ref('')
const myAvatar = ref('')
const myFirstName = ref('')
const myLastName = ref('')
const myFollowers = ref([])
const newContent = ref('')
const privacy = ref('public')
const allowedUsers = ref([])
const selectedImage = ref(null)
const imagePreview = ref('')
const fileRef = ref(null)
const states = reactive({})

// Max image size enforced both client and server side (5 MB)
const MAX_IMG_SIZE = 5 * 1024 * 1024

const myInitials = computed(() => initials(myFirstName.value, myLastName.value))

const COLORS = ['#f5a623','#3ecf78','#3b8aff','#e04040','#a78bfa','#f472b6']
function avatarColor(id) { let n=0; for(let c of (id||'')) n+=c.charCodeAt(0); return COLORS[n%COLORS.length] }
function initials(f,l) { return `${(f||'?')[0]}${(l||'?')[0]}`.toUpperCase() }
function privLabel(p) { return p==='public'?'Public':p==='followers'?'Followers':'Private' }
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

function ensureState(id) {
  if (!states[id]) states[id] = { open:false, loading:false, loadError:false, comments:[], text:'', sending:false, commentError:'', commentImage:null, commentImagePreview:'' }
}

function pickImage(e) {
  const f = e.target.files?.[0]; if (!f) return
  if (f.size > MAX_IMG_SIZE) {
    postError.value = 'Image is too large. Maximum size is 5 MB.'
    if (fileRef.value) fileRef.value.value = ''
    return
  }
  postError.value = ''
  selectedImage.value = f; imagePreview.value = URL.createObjectURL(f)
}
function removeImage() {
  selectedImage.value = null
  if (imagePreview.value) URL.revokeObjectURL(imagePreview.value)
  imagePreview.value = ''
  if (fileRef.value) fileRef.value.value = ''
}

function pickCommentImage(postId, e) {
  const f = e.target.files?.[0]; if (!f) return
  ensureState(postId)
  if (f.size > MAX_IMG_SIZE) {
    states[postId].commentError = 'Image is too large. Maximum size is 5 MB.'
    return
  }
  states[postId].commentError = ''
  states[postId].commentImage = f
  states[postId].commentImagePreview = URL.createObjectURL(f)
}
function removeCommentImage(postId) {
  if (states[postId].commentImagePreview) URL.revokeObjectURL(states[postId].commentImagePreview)
  states[postId].commentImage = null
  states[postId].commentImagePreview = ''
}

async function loadFeed() {
  loading.value = true; fetchError.value = ''
  try {
    const data = await apiFetch('/feed')
    posts.value = Array.isArray(data) ? data : []
    posts.value.forEach(p => {
      ensureState(p.id)
      states[p.id]._pendingCommentCount = p.comment_count || 0
    })
  } catch (e) { fetchError.value = e.message } finally { loading.value = false }
}

async function loadMyFollowers() {
  try {
    const d = await apiFetch('/followers?user_id=' + myId.value)
    myFollowers.value = Array.isArray(d) ? d : []
  } catch {}
}

async function createPost() {
  posting.value = true; postError.value = ''
  try {
    let imagePath = ''
    if (selectedImage.value) {
      const fd = new FormData(); fd.append('image', selectedImage.value)
      imagePath = await apiFetch('/upload/post-image', { method:'POST', body:fd })
    }
    const body = { content: newContent.value, image: imagePath, privacy: privacy.value }
    if (privacy.value === 'private') body.allowed_user_ids = allowedUsers.value
    await apiFetch('/posts/create', { method:'POST', body:JSON.stringify(body) })
    newContent.value = ''; removeImage(); allowedUsers.value = []; await loadFeed()
  } catch (e) { postError.value = e.message } finally { posting.value = false }
}

async function react(post, type) {
  const prev = post.my_reaction; const removing = prev === type
  post.my_reaction = removing ? '' : type
  if (prev==='like') post.likes = Math.max(0, post.likes-1)
  if (prev==='dislike') post.dislikes = Math.max(0, post.dislikes-1)
  if (!removing) { if (type==='like') post.likes++; else post.dislikes++ }
  try {
    const res = await apiFetch('/posts/react', { method:'POST', body:JSON.stringify({ post_id:post.id, type:removing?'':type }) })
    post.likes=res.likes; post.dislikes=res.dislikes; post.my_reaction=res.my_reaction
  } catch { post.my_reaction=prev }
}

async function toggleComments(post) {
  ensureState(post.id); const st = states[post.id]
  st.open = !st.open
  if (st.open) {
    const needsRefresh = st.comments.length === 0 ||
      st._pendingCommentCount > st.comments.length
    if (needsRefresh) {
      st.loading = true; st.loadError = false
      try { const d = await apiFetch(`/posts/comments?post_id=${post.id}`); st.comments = Array.isArray(d)?d:[] }
      catch { st.loadError = true } finally { st.loading = false }
    }
  }
}

async function addComment(post) {
  const st = states[post.id]
  if ((!st.text?.trim() && !st.commentImage) || st.sending) return
  st.sending = true; st.commentError = ''
  try {
    let imagePath = ''
    if (st.commentImage) {
      const fd = new FormData(); fd.append('image', st.commentImage)
      imagePath = await apiFetch('/upload/comment-image', { method:'POST', body:fd })
      removeCommentImage(post.id)
    }
    await apiFetch('/posts/comments/create', { method:'POST', body:JSON.stringify({ post_id:post.id, content:st.text.trim(), image:imagePath }) })
    st.text = ''
    const d = await apiFetch(`/posts/comments?post_id=${post.id}`); st.comments = Array.isArray(d)?d:[]
    post.comment_count = st.comments.length
    states[post.id]._pendingCommentCount = st.comments.length
  } catch (e) { st.commentError = e.message || 'Failed to post comment.' }
  finally { st.sending = false }
}

async function likeComment(post, c) {
  const prev = c.i_liked; c.i_liked = !prev; c.likes = prev ? Math.max(0,c.likes-1) : c.likes+1
  try {
    const res = await apiFetch('/posts/comments/like', { method:'POST', body:JSON.stringify({ comment_id:c.id }) })
    c.likes=res.likes; c.i_liked=res.i_liked
  } catch { c.i_liked=prev }
}

async function deletePost(post) {
  if (!confirm('Delete this post?')) return
  try {
    await apiFetch(`/posts/delete?post_id=${post.id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== post.id)
  } catch (e) { postError.value = e.message || 'Failed to delete post.' }
}

async function deleteComment(post, comment) {
  if (!confirm('Delete this comment?')) return
  const st = states[post.id]
  try {
    await apiFetch(`/posts/comments/delete?comment_id=${comment.id}`, { method: 'DELETE' })
    st.comments = st.comments.filter(c => c.id !== comment.id)
    post.comment_count = Math.max(0, (post.comment_count || 0) - 1)
    states[post.id]._pendingCommentCount = st.comments.length
  } catch (e) { st.commentError = e.message || 'Failed to delete comment.' }
}

// ── Feed WebSocket: live reaction & comment-count updates ───────────────────
let feedSocket = null
let feedRetry = null
let _feedDestroyed = false

function connectFeedWS() {
  if (_feedDestroyed) return
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  feedSocket = new WebSocket(`${proto}//${location.host}/ws/feed`)
  feedSocket.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      if (msg.event === 'post_updated') {
        const post = posts.value.find(p => p.id === msg.post_id)
        if (post) {
          if (msg.likes    !== undefined) post.likes    = msg.likes
          if (msg.dislikes !== undefined) post.dislikes = msg.dislikes
          if (msg.comment_count !== undefined) {
            // Only update when comments panel is closed (to avoid jarring reloads)
            if (!states[post.id]?.open) {
              ensureState(post.id)
              // We only have the count; update when user opens comments next time
              states[post.id]._pendingCommentCount = msg.comment_count
            }
          }
        }
      }
    } catch {}
  }
  feedSocket.onclose = () => {
    if (!_feedDestroyed) feedRetry = setTimeout(connectFeedWS, 5000)
  }
}

onMounted(async () => {
  if (isLoggedIn.value) {
    try {
      const me = await apiFetch('/me')
      myId.value=me.id; myFirstName.value=me.first_name; myLastName.value=me.last_name; myAvatar.value=me.avatar||''
    } catch { router.push('/login'); return }
    await Promise.all([loadFeed(), loadMyFollowers()])
    connectFeedWS()
  } else {
    await loadFeed()
  }
})

onUnmounted(() => {
  _feedDestroyed = true
  clearTimeout(feedRetry)
  if (feedSocket) feedSocket.close(1000)
})
</script>
