<template>
  <div>
    <Navbar />
    <div class="page-layout" style="max-width:900px">
      <p v-if="loading" style="color:var(--text-3);text-align:center;padding:40px">Loading…</p>

      <template v-if="!loading && user">
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;align-items:start">

          <!-- LEFT: My profile card + posts -->
          <div style="display:flex;flex-direction:column;gap:16px">

            <!-- Profile hero -->
            <div class="card" style="padding:24px;text-align:center">
              <div class="avatar-ring" style="width:80px;height:80px;font-size:1.5rem;margin:0 auto 14px;cursor:pointer" :style="{ background: avatarColor(user.id) }" @click="triggerAvatarUpload">
                <img v-if="user.avatar" :src="BACKEND_URL + user.avatar" />
                <span v-else>{{ initials }}</span>
              </div>
              <input type="file" accept=".jpg,.jpeg,.png,.gif" style="display:none" ref="avatarInput" @change="uploadAvatar" />
              <div style="font-family:'Syne',sans-serif;font-weight:700;font-size:1.1rem;margin-bottom:2px">{{ user.first_name }} {{ user.last_name }}</div>
              <div v-if="user.nickname" style="color:var(--text-3);font-size:0.83rem;margin-bottom:4px">@{{ user.nickname }}</div>
              <div v-if="user.about_me" style="color:var(--text-2);font-size:0.85rem;line-height:1.5;margin-bottom:10px">{{ user.about_me }}</div>
              <div style="font-size:0.78rem;color:var(--text-3);margin-bottom:14px;display:flex;align-items:center;justify-content:center;gap:4px">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
                {{ user.email }}
              </div>

              <!-- Stats -->
              <div style="display:flex;border-top:1px solid var(--border);padding-top:14px;margin-top:4px">
                <div style="flex:1;text-align:center;border-right:1px solid var(--border)">
                  <div style="font-family:'Syne',sans-serif;font-weight:800;font-size:1.2rem">{{ posts.length }}</div>
                  <div style="font-size:0.72rem;color:var(--text-3);text-transform:uppercase">Posts</div>
                </div>
                <div style="flex:1;text-align:center;border-right:1px solid var(--border);cursor:pointer" @click="showFollowers=true">
                  <div style="font-family:'Syne',sans-serif;font-weight:800;font-size:1.2rem">{{ followersCount }}</div>
                  <div style="font-size:0.72rem;color:var(--text-3);text-transform:uppercase">Followers</div>
                </div>
                <div style="flex:1;text-align:center;cursor:pointer" @click="showFollowing=true">
                  <div style="font-family:'Syne',sans-serif;font-weight:800;font-size:1.2rem">{{ followingCount }}</div>
                  <div style="font-size:0.72rem;color:var(--text-3);text-transform:uppercase">Following</div>
                </div>
              </div>
            </div>

            <div class="card" style="padding:20px">
              <h3 style="font-size:0.95rem;margin-bottom:14px;padding-bottom:10px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:6px">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path></svg>
                Edit Info
              </h3>
              <div class="form-field"><label>Nickname</label><input class="input" v-model="editNickname" placeholder="@handle" /></div>
              <div class="form-field"><label>About me</label><textarea class="input" v-model="editAboutMe" rows="2" placeholder="Tell the world about yourself…"></textarea></div>

              <!-- Profile visibility toggles -->
              <div style="margin-bottom:14px">
                <div style="font-size:0.78rem;font-weight:600;color:var(--text-2);text-transform:uppercase;letter-spacing:0.05em;margin-bottom:10px">Show on profile</div>
                <div style="display:flex;flex-direction:column;gap:8px">
                  <label style="display:flex;align-items:center;justify-content:space-between;padding:10px 12px;border-radius:8px;border:1px solid var(--border);background:var(--surface2);cursor:pointer">
                    <div style="display:flex;align-items:center;gap:8px">
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
                      <span style="font-size:0.85rem;font-weight:500">Email address</span>
                    </div>
                    <div class="toggle-switch" :class="{ active: editShowEmail }" @click="editShowEmail = !editShowEmail"><div class="toggle-thumb"></div></div>
                  </label>
                  <label style="display:flex;align-items:center;justify-content:space-between;padding:10px 12px;border-radius:8px;border:1px solid var(--border);background:var(--surface2);cursor:pointer">
                    <div style="display:flex;align-items:center;gap:8px">
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                      <span style="font-size:0.85rem;font-weight:500">Date of birth</span>
                    </div>
                    <div class="toggle-switch" :class="{ active: editShowDOB }" @click="editShowDOB = !editShowDOB"><div class="toggle-thumb"></div></div>
                  </label>
                </div>
              </div>

              <button class="btn btn-primary btn-sm" :disabled="savingInfo" @click="saveInfo">{{ savingInfo ? 'Saving…' : 'Save' }}</button>
              <div v-if="infoOk" class="success-msg" style="margin-top:8px;display:flex;align-items:center;gap:4px"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> Saved!</div>
              <div v-if="infoErr" class="error-msg" style="margin-top:8px">{{ infoErr }}</div>
            </div>

            <!-- My Posts -->
            <div>
              <div class="section-label" style="margin-bottom:10px">My Posts</div>
              <div v-if="postsLoading" class="card" style="padding:20px;text-align:center;color:var(--text-3)">Loading…</div>
              <div v-else-if="!posts.length" class="card empty-state" style="padding:36px">
                <div class="empty-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
                </div>
                <h3>No posts yet</h3>
              </div>
              <article class="post-card" style="margin-bottom:10px" v-for="post in posts" :key="post.id">
                <div class="post-header">
                  <div class="avatar-ring post-avatar" :style="{ background: avatarColor(user.id) }">
                    <img v-if="user.avatar" :src="BACKEND_URL + user.avatar" />
                    <span v-else>{{ initials }}</span>
                  </div>
                  <div class="post-meta">
                    <span class="post-author">You</span>
                    <span class="post-time">{{ relTime(post.created_at) }}</span>
                  </div>
                  <span class="privacy-tag" style="display:inline-flex;align-items:center;gap:4px">
                    <svg v-if="post.privacy === 'public'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                    <svg v-else-if="post.privacy === 'followers'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                    {{ post.privacy === 'public' ? 'Public' : post.privacy === 'followers' ? 'Followers' : 'Private' }}
                  </span>
                </div>
                <p class="post-body" v-if="post.content">{{ post.content }}</p>
                <img v-if="post.image" :src="BACKEND_URL + post.image" class="post-img" />
              </article>
            </div>
          </div>

          <!-- RIGHT: Settings -->
          <div style="display:flex;flex-direction:column;gap:16px">

            <!-- Avatar upload -->
            <div class="card" style="padding:20px">
              <h3 style="font-size:0.95rem;margin-bottom:14px;padding-bottom:10px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:6px">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                Avatar
              </h3>
              <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap">
                <label style="cursor:pointer">
                  <input type="file" accept=".jpg,.jpeg,.png,.gif" style="display:none" @change="uploadAvatar" />
                  <span class="btn btn-ghost btn-sm" style="display:inline-flex;align-items:center;gap:4px">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
                    Choose photo
                  </span>
                </label>
              </div>
              <div v-if="avatarOk" class="success-msg" style="margin-top:8px;display:flex;align-items:center;gap:4px"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> Avatar updated!</div>
              <div v-if="avatarErr" class="error-msg" style="margin-top:8px">{{ avatarErr }}</div>
            </div>

            <!-- Password -->
            <div class="card" style="padding:20px">
              <h3 style="font-size:0.95rem;margin-bottom:14px;padding-bottom:10px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:6px">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                Change Password
              </h3>
              <div class="form-field"><label>Current password</label><input class="input" type="password" v-model="pwCurrent" /></div>
              <div class="form-field"><label>New password</label><input class="input" type="password" v-model="pwNew" /></div>
              <div class="form-field"><label>Confirm new password</label><input class="input" type="password" v-model="pwConfirm" /></div>
              <div v-if="pwConfirm && pwNew !== pwConfirm" style="color:var(--danger);font-size:0.8rem;margin-bottom:10px">Passwords don't match</div>
              <button class="btn btn-primary btn-sm" :disabled="!pwCurrent||!pwNew||pwNew!==pwConfirm||pwNew.length<8||pwChanging" @click="changePassword">
                {{ pwChanging ? 'Saving…' : 'Update password' }}
              </button>
              <div v-if="pwOk" class="success-msg" style="margin-top:8px;display:flex;align-items:center;gap:4px"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> Password updated!</div>
              <div v-if="pwErr" class="error-msg" style="margin-top:8px">{{ pwErr }}</div>
            </div>

            <div class="card" style="padding:20px">
              <h3 style="font-size:0.95rem;margin-bottom:14px;padding-bottom:10px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:6px">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                Privacy
              </h3>
              <div style="display:flex;flex-direction:column;gap:10px;margin-bottom:14px">
                <label v-for="opt in [{v:'public',t:'Public',d:'Anyone can view your profile'},{v:'private',t:'Private',d:'Only followers see your content'}]" :key="opt.v"
                  :style="{ padding:'12px 14px', borderRadius:'10px', border:'2px solid', borderColor: privacyType===opt.v ? 'var(--accent)' : 'var(--border)', background: privacyType===opt.v ? 'var(--accent-bg)' : 'var(--surface2)', cursor:'pointer', display:'flex', alignItems:'center', gap:'10px' }">
                  <input type="radio" v-model="privacyType" :value="opt.v" style="display:none" />
                  <span style="font-size:1.3rem;display:flex;align-items:center;color:var(--text)">
                    <svg v-if="opt.v === 'public'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                  </span>
                  <div>
                    <div style="font-weight:600;font-size:0.88rem">{{ opt.t }}</div>
                    <div style="font-size:0.76rem;color:var(--text-2)">{{ opt.d }}</div>
                  </div>
                </label>
              </div>
              <button class="btn btn-primary btn-sm" :disabled="privacyType===user.profile_type||privSaving" @click="savePrivacy">
                {{ privSaving ? 'Saving…' : 'Save privacy' }}
              </button>
              <div v-if="privOk" class="success-msg" style="margin-top:8px;display:flex;align-items:center;gap:4px"><svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg> Updated!</div>
              <div v-if="privErr" class="error-msg" style="margin-top:8px">{{ privErr }}</div>
            </div>

          </div>
        </div>
      </template>
    </div>

    <!-- Followers modal -->
    <div v-if="showFollowers" class="modal-backdrop" @click.self="showFollowers=false">
      <div class="modal-box">
        <div class="modal-header"><h3>Followers</h3><button class="modal-close" @click="showFollowers=false"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg></button></div>
        <div v-if="!followersList.length" style="color:var(--text-3);text-align:center;padding:24px">No followers yet</div>
        <div v-for="u in followersList" :key="u.user_id" style="display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid var(--border)">
          <router-link :to="'/profile/'+u.user_id" @click="showFollowers=false" style="display:flex;align-items:center;gap:10px;color:var(--text);text-decoration:none;flex:1">
            <div class="avatar-ring" style="width:34px;height:34px;font-size:0.72rem" :style="{ background: avatarColor(u.user_id) }">{{ initials2(u.first_name,u.last_name) }}</div>
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
        <div v-for="u in followingList" :key="u.id" style="display:flex;align-items:center;gap:10px;padding:8px 0;border-bottom:1px solid var(--border)">
          <router-link :to="'/profile/'+u.id" @click="showFollowing=false" style="display:flex;align-items:center;gap:10px;color:var(--text);text-decoration:none;flex:1">
            <div class="avatar-ring" style="width:34px;height:34px;font-size:0.72rem" :style="{ background: avatarColor(u.id) }">{{ initials2(u.first_name,u.last_name) }}</div>
            <span style="font-weight:600;font-size:0.9rem">{{ u.first_name }} {{ u.last_name }}</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch, BACKEND_URL } from '../services/api'

const router = useRouter()
const user = ref(null)
const loading = ref(true)
const posts = ref([])
const postsLoading = ref(true)
const followersList = ref([])
const followingList = ref([])
const followersCount = ref(0)
const followingCount = ref(0)
const showFollowers = ref(false)
const showFollowing = ref(false)
const avatarInput = ref(null)

// Edit info
const editNickname = ref('')
const editAboutMe = ref('')
const editShowEmail = ref(false)
const editShowDOB = ref(false)
const savingInfo = ref(false)
const infoOk = ref(false)
const infoErr = ref('')

// Avatar
const avatarOk = ref(false)
const avatarErr = ref('')

// Password
const pwCurrent = ref(''); const pwNew = ref(''); const pwConfirm = ref('')
const pwChanging = ref(false); const pwOk = ref(false); const pwErr = ref('')

// Privacy
const privacyType = ref('public')
const privSaving = ref(false); const privOk = ref(false); const privErr = ref('')

const initials = computed(() => user.value ? `${user.value.first_name[0]}${user.value.last_name[0]}`.toUpperCase() : '?')
function initials2(f,l) { return `${(f||'?')[0]}${(l||'?')[0]}`.toUpperCase() }

const COLORS = ['#f5a623','#3ecf78','#3b8aff','#e04040','#a78bfa','#f472b6']
function avatarColor(id) { let n=0; for(let c of (id||'')) n+=c.charCodeAt(0); return COLORS[n%COLORS.length] }

function privLabel(p) { return p==='public'?'Public':p==='followers'?'Followers':'Private' }
function relTime(s) {
  if (!s) return ''
  const normalized = s.includes('T') ? s : s.replace(' ', 'T') + 'Z'
  const ts = new Date(normalized).getTime()
  if (isNaN(ts)) return ''
  const d = Date.now() - ts
  if (d < 0) return 'just now'
  const sec=Math.floor(d/1000); if(sec<60) return 'just now'
  const m=Math.floor(sec/60); if(m<60) return m+'m ago'
  const h=Math.floor(m/60); if(h<24) return h+'h ago'
  const day=Math.floor(h/24); if(day<7) return day+'d ago'
  return new Date(normalized).toLocaleDateString()
}

function triggerAvatarUpload() { avatarInput.value?.click() }

async function uploadAvatar(e) {
  const file = e.target.files?.[0]; if (!file) return
  avatarErr.value=''; avatarOk.value=false
  try {
    const fd = new FormData(); fd.append('avatar', file)
    const path = await apiFetch('/upload/avatar', { method:'POST', body:fd })
    user.value.avatar = path; avatarOk.value=true
    setTimeout(()=>avatarOk.value=false, 3000)
  } catch(e) { avatarErr.value=e.message }
}

async function saveInfo() {
  savingInfo.value=true; infoErr.value=''; infoOk.value=false
  try {
    await apiFetch('/profile/update', { method:'POST', body:JSON.stringify({
      nickname: editNickname.value,
      about_me: editAboutMe.value,
      show_email: editShowEmail.value,
      show_dob: editShowDOB.value,
    })})
    user.value.nickname = editNickname.value
    user.value.about_me = editAboutMe.value
    user.value.show_email = editShowEmail.value
    user.value.show_dob = editShowDOB.value
    infoOk.value=true; setTimeout(()=>infoOk.value=false,3000)
  } catch(e) { infoErr.value=e.message || 'Failed to update info.' } finally { savingInfo.value=false }
}

async function changePassword() {
  pwChanging.value=true; pwErr.value=''; pwOk.value=false
  try {
    await apiFetch('/profile/password', { method:'POST', body:JSON.stringify({ current_password:pwCurrent.value, new_password:pwNew.value }) })
    pwCurrent.value=''; pwNew.value=''; pwConfirm.value=''
    pwOk.value=true; setTimeout(()=>pwOk.value=false,3000)
  } catch(e) { pwErr.value=e.message } finally { pwChanging.value=false }
}

async function savePrivacy() {
  if (!confirm(`Switch to ${privacyType.value} profile?`)) { privacyType.value=user.value.profile_type; return }
  privSaving.value=true; privErr.value=''; privOk.value=false
  try {
    await apiFetch('/profile/privacy', { method:'POST', body:JSON.stringify({ profile_type:privacyType.value }) })
    user.value.profile_type=privacyType.value
    privOk.value=true; setTimeout(()=>privOk.value=false,3000)
  } catch(e) { privErr.value=e.message } finally { privSaving.value=false }
}

onMounted(async () => {
  try {
    user.value = await apiFetch('/me')
    privacyType.value = user.value.profile_type||'public'
    editNickname.value = user.value.nickname||''
    editAboutMe.value = user.value.about_me||''
    editShowEmail.value = !!user.value.show_email
    editShowDOB.value = !!user.value.show_dob
  } catch { router.push('/login'); return } finally { loading.value=false }

  try {
    const d = await apiFetch(`/profile/posts?user_id=${user.value.id}`)
    posts.value = Array.isArray(d)?d:[]
  } catch {} finally { postsLoading.value=false }

  try { followersList.value = await apiFetch(`/followers?user_id=${user.value.id}`); followersCount.value=followersList.value.length } catch {}
  try { followingList.value = await apiFetch(`/following?user_id=${user.value.id}`); followingCount.value=followingList.value.length } catch {}
})
</script>

<style scoped>
.toggle-switch {
  width: 38px; height: 20px;
  background: var(--border);
  border-radius: 10px;
  position: relative;
  transition: background 0.2s;
  flex-shrink: 0;
}
.toggle-switch.active { background: var(--accent); }
.toggle-thumb {
  position: absolute;
  top: 2px; left: 2px;
  width: 16px; height: 16px;
  border-radius: 50%;
  background: #fff;
  transition: left 0.2s;
  box-shadow: 0 1px 3px rgba(0,0,0,0.25);
}
.toggle-switch.active .toggle-thumb { left: 20px; }
</style>
