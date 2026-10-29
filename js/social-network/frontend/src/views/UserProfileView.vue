<template>
  <div>
    <Navbar />
    <div class="page-layout" style="max-width:780px">

      <!-- Skeleton -->
      <template v-if="loading">
        <div class="card" style="padding:28px">
          <div style="display:flex;gap:18px;align-items:flex-start;flex-wrap:wrap">
            <div class="skeleton" style="width:90px;height:90px;border-radius:50%;flex-shrink:0"></div>
            <div style="flex:1;display:flex;flex-direction:column;gap:10px;padding-top:8px">
              <div class="skeleton" style="height:16px;width:45%"></div>
              <div class="skeleton" style="height:12px;width:28%"></div>
              <div class="skeleton" style="height:12px;width:70%"></div>
            </div>
          </div>
        </div>
      </template>

      <!-- Error -->
      <div v-else-if="fetchError" class="card empty-state">
        <p style="display:inline-flex;align-items:center;gap:6px;justify-content:center;color:var(--danger)">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
          {{ fetchError }}
        </p>
        <button class="btn btn-secondary btn-sm" style="margin-top:12px" @click="load">Retry</button>
      </div>

      <template v-else-if="profile">
        <!-- Profile hero -->
        <div class="card" style="padding:28px 24px">
          <div style="display:flex;gap:20px;align-items:flex-start;flex-wrap:wrap">
            <div class="avatar-ring" style="width:88px;height:88px;font-size:1.6rem;flex-shrink:0" :style="{ background: avatarColor(profile.id) }">
              <img v-if="profile.avatar" :src="BACKEND_URL + profile.avatar" />
              <span v-else>{{ initials(profile.first_name, profile.last_name) }}</span>
            </div>

            <div style="flex:1;min-width:0">
              <div style="display:flex;align-items:flex-start;gap:12px;flex-wrap:wrap;justify-content:space-between">
                <div>
                  <h1 style="font-size:1.35rem;margin-bottom:2px">{{ profile.first_name }} {{ profile.last_name }}</h1>
                  <p v-if="profile.nickname" style="color:var(--text-3);font-size:0.85rem;margin-bottom:6px">@{{ profile.nickname }}</p>
                  <p v-if="profile.about_me" style="color:var(--text-2);font-size:0.88rem;line-height:1.5;margin-bottom:8px">{{ profile.about_me }}</p>

                  <!-- Conditionally visible contact info -->
                  <div v-if="profile.email" style="display:inline-flex;align-items:center;gap:5px;font-size:0.8rem;color:var(--text-3);margin-bottom:4px">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
                    {{ profile.email }}
                  </div>
                  <div v-if="profile.date_of_birth" style="display:inline-flex;align-items:center;gap:5px;font-size:0.8rem;color:var(--text-3);margin-bottom:4px;margin-left:4px">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                    {{ profile.date_of_birth }}
                  </div>

                  <span v-if="profile.profile_type==='private'" style="font-size:0.78rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                    Private account
                  </span>
                </div>

                <!-- Actions -->
                <div style="display:flex;gap:8px;align-items:center">
                  <router-link v-if="profile.is_owner" to="/profile" class="btn btn-ghost btn-sm">Edit Profile</router-link>
                  <template v-else>
                    <button v-if="profile.follow_status==='none'" class="btn btn-primary btn-sm" :disabled="actionLoading" @click="follow">
                      {{ actionLoading ? '…' : 'Follow' }}
                    </button>
                    <button v-else-if="profile.follow_status==='pending'" class="btn btn-ghost btn-sm" :disabled="actionLoading" @click="cancelFollow">
                      <svg v-if="!actionLoading" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                      {{ actionLoading ? '…' : 'Requested' }}
                    </button>
                    <div v-else-if="profile.follow_status==='accepted'" style="display:flex;gap:8px">
                      <button class="btn btn-danger btn-sm" :disabled="actionLoading" @click="unfollow">
                        {{ actionLoading ? '…' : 'Unfollow' }}
                      </button>
                      <button class="btn btn-secondary btn-sm" style="display:inline-flex;align-items:center;gap:6px" @click="goToDm">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
                        Message
                      </button>
                    </div>
                  </template>
                </div>
              </div>

              <!-- Stats -->
              <div style="display:flex;gap:0;margin-top:18px;border-top:1px solid var(--border);padding-top:16px">
                <div v-for="stat in [{n:profile.posts_count,l:'Posts',click:null},{n:profile.followers_count,l:'Followers',click:()=>showFollowers=true},{n:profile.following_count,l:'Following',click:()=>showFollowing=true}]" :key="stat.l"
                  style="flex:1;text-align:center;padding:4px 0"
                  :style="{ borderRight: stat.l!=='Following' ? '1px solid var(--border)' : 'none', cursor: stat.click ? 'pointer' : 'default' }"
                  @click="stat.click && stat.click()">
                  <div style="font-family:'Syne',sans-serif;font-weight:800;font-size:1.3rem;color:var(--text)">{{ stat.n }}</div>
                  <div style="font-size:0.72rem;color:var(--text-3);text-transform:uppercase;letter-spacing:0.06em">{{ stat.l }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Private wall -->
        <div v-if="!profile.can_view_posts && !profile.is_owner" class="card empty-state" style="padding:60px 24px">
          <div class="empty-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
          </div>
          <h3>This account is private</h3>
          <p>Follow to see their posts</p>
          <button v-if="profile.follow_status==='none'" class="btn btn-primary" style="margin-top:18px" :disabled="actionLoading" @click="follow">
            {{ actionLoading ? '…' : 'Send Follow Request' }}
          </button>
          <p v-else-if="profile.follow_status==='pending'" style="margin-top:14px;color:var(--text-3);font-size:0.88rem;display:inline-flex;align-items:center;gap:6px">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            Request pending…
          </p>
        </div>

        <!-- Posts -->
        <div v-else>
          <div class="section-label">Posts</div>

          <div v-if="postsLoading">
            <div class="card" style="padding:20px;margin-bottom:12px" v-for="n in 2" :key="n">
              <div style="display:flex;gap:12px;margin-bottom:14px">
                <div class="skeleton" style="width:40px;height:40px;border-radius:50%;flex-shrink:0"></div>
                <div style="flex:1;display:flex;flex-direction:column;gap:6px;padding-top:4px">
                  <div class="skeleton" style="height:12px;width:38%"></div>
                  <div class="skeleton" style="height:10px;width:22%"></div>
                </div>
              </div>
              <div class="skeleton" style="height:13px;width:100%;margin-bottom:8px"></div>
              <div class="skeleton" style="height:13px;width:60%"></div>
            </div>
          </div>

          <div v-else-if="posts.length===0" class="card empty-state" style="padding:48px">
            <div class="empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            </div>
            <h3>No posts yet</h3>
          </div>

          <article v-else class="post-card" style="margin-bottom:12px" v-for="post in posts" :key="post.id">
            <div class="post-header">
              <div class="avatar-ring post-avatar" :style="{ background: avatarColor(post.user_id) }">
                <img v-if="post.avatar" :src="BACKEND_URL + post.avatar" />
                <span v-else>{{ initials(post.first_name, post.last_name) }}</span>
              </div>
              <div class="post-meta">
                <span class="post-author">{{ post.first_name }} {{ post.last_name }}</span>
                <span class="post-time">{{ relTime(post.created_at) }}</span>
              </div>
              <span class="privacy-tag" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="post.privacy === 'public'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                <svg v-else-if="post.privacy === 'followers'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                {{ post.privacy === 'public' ? 'Public' : post.privacy === 'followers' ? 'Followers' : 'Private' }}
              </span>
              <button v-if="profile.is_owner" class="btn btn-danger btn-sm" style="margin-left:auto;display:inline-flex;align-items:center;justify-content:center;padding:6px 8px" @click="deletePost(post)">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
              </button>
            </div>
            <p class="post-body" v-if="post.content">{{ post.content }}</p>
            <img v-if="post.image" :src="BACKEND_URL + post.image" class="post-img" />
          </article>
        </div>
      </template>
    </div>

    <!-- Followers modal -->
    <div v-if="showFollowers" class="modal-backdrop" @click.self="showFollowers=false">
      <div class="modal-box">
        <div class="modal-header"><h3>Followers</h3><button class="modal-close" @click="showFollowers=false"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg></button></div>
        <div v-if="!followersList.length" style="color:var(--text-3);text-align:center;padding:24px">No followers yet</div>
        <div v-for="u in followersList" :key="u.user_id" style="display:flex;align-items:center;gap:10px;padding:8px 4px;border-bottom:1px solid var(--border)">
          <router-link :to="'/profile/'+u.user_id" @click="showFollowers=false" style="display:flex;align-items:center;gap:10px;color:var(--text);text-decoration:none;flex:1">
            <div class="avatar-ring" style="width:36px;height:36px;font-size:0.75rem;flex-shrink:0" :style="{ background: avatarColor(u.user_id) }">{{ initials(u.first_name,u.last_name) }}</div>
            <span style="font-weight:600;font-size:0.9rem">{{ u.first_name }} {{ u.last_name }}</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Following modal -->
    <div v-if="showFollowing" class="modal-backdrop" @click.self="showFollowing=false">
      <div class="modal-box">
        <div class="modal-header"><h3>Following</h3><button class="modal-close" @click="showFollowing=false"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg></button></div>
        <div v-if="!followingList.length" style="color:var(--text-3);text-align:center;padding:24px">Not following anyone</div>
        <div v-for="u in followingList" :key="u.id" style="display:flex;align-items:center;gap:10px;padding:8px 4px;border-bottom:1px solid var(--border)">
          <router-link :to="'/profile/'+u.id" @click="showFollowing=false" style="display:flex;align-items:center;gap:10px;color:var(--text);text-decoration:none;flex:1">
            <div class="avatar-ring" style="width:36px;height:36px;font-size:0.75rem;flex-shrink:0" :style="{ background: avatarColor(u.id) }">{{ initials(u.first_name,u.last_name) }}</div>
            <span style="font-weight:600;font-size:0.9rem">{{ u.first_name }} {{ u.last_name }}</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'

const route = useRoute(); const router = useRouter()
const profile = ref(null); const posts = ref([]); const loading = ref(true)
const postsLoading = ref(false); const fetchError = ref(''); const actionLoading = ref(false)
const showFollowers = ref(false); const showFollowing = ref(false)
const followersList = ref([]); const followingList = ref([])

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

async function load() {
  loading.value=true; fetchError.value=''
  const uid = route.params.id
  try { profile.value = await apiFetch(`/profile?user_id=${uid}`) }
  catch(e) { fetchError.value=e.message; loading.value=false; return }
  loading.value=false
  if (profile.value.can_view_posts || profile.value.is_owner) {
    await loadPosts()
    try { followersList.value = await apiFetch(`/followers?user_id=${uid}`) } catch {}
    try { followingList.value = await apiFetch(`/following?user_id=${uid}`) } catch {}
  }
}

async function loadPosts() {
  postsLoading.value=true
  try { posts.value = await apiFetch(`/profile/posts?user_id=${route.params.id}`) }
  catch { posts.value=[] } finally { postsLoading.value=false }
}

async function follow() {
  if (actionLoading.value) return
  actionLoading.value=true
  try {
    await apiFetch('/follow', { method:'POST', body:JSON.stringify({ target_user_id:route.params.id }) })
    if (profile.value.profile_type==='private') { profile.value.follow_status='pending' }
    else { profile.value.follow_status='accepted'; profile.value.is_follower=true; profile.value.can_view_posts=true; profile.value.followers_count++; await loadPosts() }
  } catch { await load() } finally { actionLoading.value=false }
}
async function unfollow() {
  if (actionLoading.value) return
  if (!confirm('Unfollow this user?')) return
  actionLoading.value=true
  try {
    await apiFetch('/unfollow', { method:'POST', body:JSON.stringify({ target_user_id:route.params.id }) })
    profile.value.follow_status='none'; profile.value.is_follower=false; profile.value.followers_count=Math.max(0,profile.value.followers_count-1)
    if (profile.value.profile_type==='private') { profile.value.can_view_posts=false; posts.value=[] }
  } catch { await load() } finally { actionLoading.value=false }
}
async function deletePost(post) {
  if (!confirm('Delete this post?')) return
  try {
    await apiFetch(`/posts/delete?post_id=${post.id}`, { method: 'DELETE' })
    posts.value = posts.value.filter(p => p.id !== post.id)
    if (profile.value) profile.value.posts_count = Math.max(0, (profile.value.posts_count || 0) - 1)
  } catch (e) { alert(e.message || 'Failed to delete post.') }
}

async function cancelFollow() {
  if (actionLoading.value) return
  if (!confirm('Cancel follow request?')) return
  actionLoading.value = true
  try {
    await apiFetch('/unfollow', { method:'POST', body:JSON.stringify({ target_user_id:route.params.id }) })
    profile.value.follow_status = 'none'
  } catch { await load() } finally { actionLoading.value = false }
}

function goToDm() {
  if (profile.value) {
    router.push('/chat?user_id=' + profile.value.id)
  }
}

watch(() => route.params.id, (newId) => {
  if (newId) { posts.value=[]; profile.value=null; showFollowers.value=false; showFollowing.value=false; load() }
})

onMounted(async () => {
  try { await apiFetch('/me') } catch { router.push('/login'); return }
  await load()
})
</script>
