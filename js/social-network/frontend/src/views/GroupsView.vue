<template>
  <div>
    <Navbar />
    <div class="page-layout two-col" style="max-width:1100px">

      <!-- Left: My groups + invites -->
      <div style="display:flex;flex-direction:column;gap:16px">

        <div class="card" style="padding:22px">
          <h3 style="font-size:1rem;margin-bottom:16px;display:flex;align-items:center;gap:6px">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            Create Group
          </h3>
          <div class="form-field">
            <label>Title</label>
            <input class="input" v-model="newTitle" placeholder="Group name" :style="{ borderColor: newTitle.length > 100 ? 'var(--danger)' : '' }" />
            <div style="text-align:right;font-size:0.7rem;margin-top:2px" :style="{ color: newTitle.length > 100 ? 'var(--danger)' : 'var(--text-3)' }">{{ newTitle.length }}/100</div>
          </div>
          <div class="form-field">
            <label>Description</label>
            <textarea class="input" v-model="newDesc" placeholder="What's this group about?" rows="2" :style="{ borderColor: newDesc.length > 1000 ? 'var(--danger)' : '' }"></textarea>
            <div style="text-align:right;font-size:0.7rem;margin-top:2px" :style="{ color: newDesc.length > 1000 ? 'var(--danger)' : 'var(--text-3)' }">{{ newDesc.length }}/1000</div>
          </div>
          <div v-if="createErr" class="error-msg" style="margin-bottom:10px">{{ createErr }}</div>
          <button class="btn btn-primary w-full" :disabled="creating || !newTitle.trim() || !newDesc.trim() || newTitle.length > 100 || newDesc.length > 1000" @click="createGroup">
            {{ creating ? 'Creating…' : 'Create Group' }}
          </button>
        </div>

        <!-- Invites -->
        <div class="card" style="padding:22px" v-if="invites.length">
          <h3 style="font-size:1rem;margin-bottom:14px;display:flex;align-items:center;gap:6px">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            Invitations <span class="badge" style="margin-left:6px">{{ invites.length }}</span>
          </h3>
          <div v-for="g in invites" :key="g.id" style="padding:12px;background:var(--surface2);border-radius:10px;margin-bottom:8px">
            <div style="font-weight:600;font-size:0.9rem;margin-bottom:2px">{{ g.title }}</div>
            <div style="font-size:0.78rem;color:var(--text-2);margin-bottom:10px">{{ g.description }}</div>
            <div style="display:flex;gap:8px">
              <button class="btn btn-primary btn-sm" @click="acceptInvite(g.id)">Accept</button>
              <button class="btn btn-ghost btn-sm" @click="declineInvite(g.id)">Decline</button>
            </div>
          </div>
        </div>

        <!-- My groups -->
        <div class="card" style="padding:22px">
          <h3 style="font-size:1rem;margin-bottom:14px;display:flex;align-items:center;gap:6px">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            My Groups
          </h3>
          <div v-if="!myGroups.length" style="color:var(--text-3);font-size:0.85rem">You haven't joined any groups yet.</div>
          <router-link v-for="g in myGroups" :key="g.id" :to="'/groups/'+g.id" style="text-decoration:none">
            <div style="padding:12px;background:var(--surface2);border-radius:10px;margin-bottom:8px;transition:background 0.15s;cursor:pointer" @mouseenter="e=>e.currentTarget.style.background='var(--surface3)'" @mouseleave="e=>e.currentTarget.style.background='var(--surface2)'">
              <div style="font-weight:600;font-size:0.9rem;color:var(--text);margin-bottom:2px">{{ g.title }}</div>
              <div style="font-size:0.76rem;color:var(--text-3)">{{ g.member_count }} members</div>
            </div>
          </router-link>
        </div>
      </div>

      <!-- Right: All groups -->
      <div>
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:14px">
          <h2 style="font-size:1.3rem">Discover Groups</h2>
        </div>
        <input class="input" v-model="searchQ" placeholder="Search groups…" style="margin-bottom:16px" />

        <div v-if="!filteredAll.length" class="card empty-state" style="padding:48px">
          <div class="empty-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
          </div>
          <h3>No groups found</h3>
        </div>

        <div class="card" style="padding:18px;margin-bottom:10px" v-for="g in filteredAll" :key="g.id">
          <div style="display:flex;align-items:flex-start;gap:14px">
            <div style="width:46px;height:46px;border-radius:12px;background:var(--accent-bg);display:flex;align-items:center;justify-content:center;color:var(--accent);flex-shrink:0">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            </div>
            <div style="flex:1;min-width:0">
              <div style="display:flex;align-items:center;justify-content:space-between;gap:10px;flex-wrap:wrap">
                <router-link :to="'/groups/'+g.id" style="font-weight:700;font-size:0.95rem;color:var(--text);text-decoration:none">{{ g.title }}</router-link>
                <span v-if="g.membership_status === 'requested'" style="font-size:0.8rem;color:var(--text-3);font-style:italic">Request Pending</span>
                <div v-else-if="g.membership_status === 'invited'" style="display:flex;gap:6px;align-items:center">
                  <span style="font-size:0.8rem;color:var(--accent);font-weight:600">Invited</span>
                  <button class="btn btn-primary btn-sm" @click="acceptInvite(g.id)">Accept</button>
                  <button class="btn btn-ghost btn-sm" @click="declineInvite(g.id)">Decline</button>
                </div>
                <button v-else-if="g.membership_status === ''" class="btn btn-ghost btn-sm" @click="requestJoin(g.id)">Join →</button>
                <router-link v-else :to="'/groups/'+g.id" class="btn btn-primary btn-sm">Open</router-link>
              </div>
              <p style="font-size:0.83rem;color:var(--text-2);margin-top:4px;line-height:1.5">{{ g.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch } from '../services/api'

const router = useRouter()
const allGroups = ref([]); const myGroups = ref([]); const invites = ref([])
const newTitle = ref(''); const newDesc = ref(''); const creating = ref(false); const createErr = ref('')
const searchQ = ref('')

const filteredAll = computed(() => allGroups.value.filter(g => g.title.toLowerCase().includes(searchQ.value.toLowerCase()) || g.description.toLowerCase().includes(searchQ.value.toLowerCase())))
function isMyGroup(id) { return myGroups.value.some(g => g.id === id) }

async function refresh() {
  try { const d=await apiFetch('/groups'); allGroups.value=Array.isArray(d)?d:[] } catch {}
  try { const d=await apiFetch('/groups/my'); myGroups.value=Array.isArray(d)?d:[] } catch {}
  try { const d=await apiFetch('/groups/invites'); invites.value=Array.isArray(d)?d:[] } catch {}
}

async function createGroup() {
  if (!newTitle.value.trim() || !newDesc.value.trim()) return
  if (newTitle.value.trim().length > 100) { createErr.value = 'Title cannot exceed 100 characters.'; return }
  if (newDesc.value.trim().length > 1000) { createErr.value = 'Description cannot exceed 1000 characters.'; return }
  creating.value=true; createErr.value=''
  try {
    await apiFetch('/groups/create', { method:'POST', body:JSON.stringify({ title:newTitle.value.trim(), description:newDesc.value.trim(), invitees:[] }) })
    newTitle.value=''; newDesc.value=''; await refresh()
  } catch(e) { createErr.value=e.message } finally { creating.value=false }
}

async function requestJoin(id) {
  try { await apiFetch('/groups/request', { method:'POST', body:JSON.stringify({ group_id:id }) }); await refresh() }
  catch(e) { alert(e.message) }
}
async function acceptInvite(id) {
  try { await apiFetch('/groups/invite/accept', { method:'POST', body:JSON.stringify({ group_id:id }) }); await refresh() }
  catch {}
}
async function declineInvite(id) {
  try { await apiFetch('/groups/invite/decline', { method:'POST', body:JSON.stringify({ group_id:id }) }); await refresh() }
  catch {}
}

onMounted(async () => {
  try { await apiFetch('/me') } catch { router.push('/login'); return }
  await refresh()
})
</script>
