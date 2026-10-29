<template>
  <div>
    <Navbar />
    <div class="page-layout" style="max-width:1080px">

      <div style="display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:12px;margin-bottom:20px">
        <div>
          <h2 style="font-size:1.6rem;margin-bottom:4px">Explore</h2>
          <p style="color:var(--text-2);font-size:0.88rem">Discover people and posts</p>
        </div>
        <div style="display:flex;gap:6px">
          <button :class="['btn btn-sm', activeTab==='people' ? 'btn-primary' : 'btn-ghost']" @click="activeTab='people'" style="display:inline-flex;align-items:center;gap:4px">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            People
          </button>
          <button :class="['btn btn-sm', activeTab==='posts' ? 'btn-primary' : 'btn-ghost']" @click="activeTab='posts'" style="display:inline-flex;align-items:center;gap:4px">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
            Posts
          </button>
        </div>
      </div>

      <!-- People tab -->
      <template v-if="activeTab==='people'">
        <input class="input" v-model="userSearch" placeholder="Search users…" style="margin-bottom:16px" />
        <div v-if="loadingUsers" style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px">
          <div class="card" style="padding:20px;text-align:center" v-for="n in 6" :key="n">
            <div class="skeleton" style="width:60px;height:60px;border-radius:50%;margin:0 auto 12px"></div>
            <div class="skeleton" style="height:12px;width:70%;margin:0 auto 8px"></div>
            <div class="skeleton" style="height:10px;width:50%;margin:0 auto"></div>
          </div>
        </div>
        <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px">
          <div class="card" style="padding:20px;text-align:center" v-for="u in filteredUsers" :key="u.id">
            <router-link :to="'/profile/'+u.id" style="text-decoration:none">
              <div class="avatar-ring" style="width:60px;height:60px;font-size:1.1rem;margin:0 auto 10px" :style="{ background: avatarColor(u.id) }">
                <img v-if="u.avatar" :src="BACKEND_URL + u.avatar" />
                <span v-else>{{ initials(u.first_name, u.last_name) }}</span>
              </div>
              <div style="font-weight:700;font-size:0.92rem;color:var(--text);margin-bottom:2px">{{ u.first_name }} {{ u.last_name }}</div>
              <div v-if="u.nickname" style="font-size:0.78rem;color:var(--text-3);margin-bottom:10px">@{{ u.nickname }}</div>
            </router-link>
            <button v-if="followStatus[u.id]==='accepted'" class="btn btn-ghost btn-sm w-full" @click="unfollow(u)">✓ Following</button>
            <button v-else-if="followStatus[u.id]==='pending'" class="btn btn-ghost btn-sm w-full" style="color:var(--text-3)">⏳ Requested</button>
            <button v-else class="btn btn-primary btn-sm w-full" @click="follow(u)">Follow</button>
          </div>
          <div v-if="!filteredUsers.length" class="card empty-state" style="grid-column:1/-1;padding:48px">
            <div class="empty-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            </div>
            <h3>{{ userSearch.trim() ? 'No users found' : 'Search for someone to follow' }}</h3>
          </div>
        </div>
      </template>

      <!-- Posts tab -->
      <template v-if="activeTab==='posts'">
        <div style="display:flex;justify-content:flex-end;gap:6px;margin-bottom:16px">
          <button v-for="s in ['new','trending']" :key="s"
            :class="['btn btn-sm', sort===s ? 'btn-primary' : 'btn-ghost']"
            @click="sort=s">
            <span style="display:inline-flex;align-items:center;gap:4px">
              <svg v-if="s === 'new'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2L2 22h20L12 2z"></path></svg>
              {{ s==='new' ? 'New' : 'Trending' }}
            </span>
          </button>
        </div>

        <div v-if="loadingPosts" class="explore-grid">
          <div class="card" style="overflow:hidden;padding:0" v-for="n in 9" :key="n">
            <div class="skeleton" style="aspect-ratio:1;width:100%"></div>
            <div style="padding:10px;display:flex;gap:8px">
              <div class="skeleton" style="width:26px;height:26px;border-radius:50%;flex-shrink:0"></div>
              <div style="flex:1;display:flex;flex-direction:column;gap:5px">
                <div class="skeleton" style="height:10px;width:60%"></div>
                <div class="skeleton" style="height:8px;width:35%"></div>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="!sorted.length" class="card empty-state">
          <div class="empty-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
          </div>
          <h3>No public posts yet</h3>
        </div>

        <div class="explore-grid" v-else>
          <div class="card" style="overflow:hidden;padding:0;cursor:pointer;transition:transform 0.18s"
            v-for="post in sorted" :key="post.id"
            @click="modal=post">
            <div style="aspect-ratio:1;overflow:hidden;background:var(--surface2)">
              <img v-if="post.image" :src="BACKEND_URL + post.image" style="width:100%;height:100%;object-fit:cover;display:block" />
              <div v-else style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;padding:14px;font-size:0.82rem;color:var(--text-2);line-height:1.5;word-break:break-word">
                {{ post.content.slice(0,120) }}
              </div>
            </div>
            <div style="display:flex;align-items:center;gap:8px;padding:10px 12px;border-top:1px solid var(--border)">
              <router-link :to="'/profile/'+post.user_id" @click.stop>
                <div class="avatar-ring" style="width:26px;height:26px;font-size:0.62rem;flex-shrink:0" :style="{ background: avatarColor(post.user_id) }">
                  <img v-if="post.avatar" :src="BACKEND_URL + post.avatar" />
                  <span v-else>{{ initials(post.first_name, post.last_name) }}</span>
                </div>
              </router-link>
              <div style="flex:1;overflow:hidden">
                <div style="font-size:0.78rem;font-weight:600;color:var(--text);white-space:nowrap;overflow:hidden;text-overflow:ellipsis">{{ post.first_name }} {{ post.last_name }}</div>
                <div style="font-size:0.68rem;color:var(--text-3)">{{ relTime(post.created_at) }}</div>
              </div>
              <div style="font-size:0.72rem;color:var(--text-3);display:inline-flex;align-items:center;gap:3px">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
                {{ post.likes }}
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Post modal -->
    <div v-if="modal" class="modal-backdrop" @click.self="modal=null">
      <div class="modal-box">
        <div class="modal-header">
          <div style="display:flex;align-items:center;gap:10px">
            <router-link :to="'/profile/'+modal.user_id" @click="modal=null">
              <div class="avatar-ring" style="width:36px;height:36px;font-size:0.78rem" :style="{ background: avatarColor(modal.user_id) }">
                <img v-if="modal.avatar" :src="BACKEND_URL + modal.avatar" />
                <span v-else>{{ initials(modal.first_name, modal.last_name) }}</span>
              </div>
            </router-link>
            <div>
              <div style="font-weight:600;font-size:0.92rem">{{ modal.first_name }} {{ modal.last_name }}</div>
              <div style="font-size:0.74rem;color:var(--text-3)">{{ relTime(modal.created_at) }}</div>
            </div>
          </div>
          <button class="modal-close" @click="modal=null">✕</button>
        </div>
        <p style="color:var(--text-2);font-size:0.93rem;line-height:1.65;white-space:pre-wrap;word-break:break-word;margin-bottom:14px">{{ modal.content }}</p>
        <img v-if="modal.image" :src="BACKEND_URL + modal.image" style="width:100%;border-radius:10px;border:1px solid var(--border);max-height:360px;object-fit:cover" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'

const router = useRouter()
const activeTab = ref('people')
const userSearch = ref('')
const allUsers = ref([])
const loadingUsers = ref(true)
const followStatus = ref({})
const myId = ref('')
const posts = ref([])
const loadingPosts = ref(true)
const sort = ref('new')
const modal = ref(null)

const sorted = computed(() => {
  if (sort.value==='trending') return [...posts.value].sort((a,b)=>(b.likes+b.dislikes)-(a.likes+a.dislikes))
  return posts.value
})

const filteredUsers = computed(() => {
  if (!userSearch.value.trim()) return []
  return allUsers.value.filter(u =>
    `${u.first_name} ${u.last_name} ${u.nickname||''}`.toLowerCase().includes(userSearch.value.toLowerCase())
  )
})

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

async function follow(user) {
  try {
    await apiFetch('/follow', { method:'POST', body:JSON.stringify({ target_user_id:user.id }) })
    const profile = await apiFetch(`/profile?user_id=${user.id}`)
    followStatus.value[user.id] = profile.follow_status
  } catch(e) { alert(e.message) }
}

async function unfollow(user) {
  if (!confirm(`Unfollow ${user.first_name}?`)) return
  try {
    await apiFetch('/unfollow', { method:'POST', body:JSON.stringify({ target_user_id:user.id }) })
    followStatus.value[user.id] = 'none'
  } catch(e) { alert(e.message) }
}

onMounted(async () => {
  try { const me=await apiFetch('/me'); myId.value=me.id }
  catch { router.push('/login'); return }

  try { const d=await apiFetch('/users'); allUsers.value=Array.isArray(d)?d:[] }
  catch {} finally { loadingUsers.value=false }

  try { const d=await apiFetch('/explore'); posts.value=Array.isArray(d)?d:[] }
  catch {} finally { loadingPosts.value=false }

  try {
    const following=await apiFetch('/following?user_id='+myId.value)
    if (Array.isArray(following)) following.forEach(u=>{ followStatus.value[u.id]='accepted' })
  } catch {}
})
</script>

<style scoped>
.explore-grid {
  display:grid;
  grid-template-columns:repeat(3,1fr);
  gap:12px;
}
@media (max-width:700px) { .explore-grid { grid-template-columns:repeat(2,1fr); } }
@media (max-width:420px) { .explore-grid { grid-template-columns:1fr; } }
</style>
