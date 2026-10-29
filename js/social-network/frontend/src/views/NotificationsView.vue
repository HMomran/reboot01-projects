<template>
  <div>
    <Navbar />
    <div class="page-layout feed-col">
      <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px">
        <h2 style="font-size:1.4rem">Notifications</h2>
        <div style="display:flex;gap:8px">
          <button v-if="notifications.length" class="btn btn-danger btn-sm" @click="clearAll" style="display:inline-flex;align-items:center;gap:4px">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
            Clear all
          </button>
          <button v-if="unread.length" class="btn btn-ghost btn-sm" :disabled="markingAll" @click="markAll" style="display:inline-flex;align-items:center;gap:4px">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            {{ markingAll ? '…' : 'Mark all read' }}
          </button>
        </div>
      </div>

      <template v-if="loading">
        <div class="card" style="padding:16px;display:flex;gap:12px;align-items:center;margin-bottom:8px" v-for="n in 5" :key="n">
          <div class="skeleton" style="width:40px;height:40px;border-radius:50%;flex-shrink:0"></div>
          <div style="flex:1;display:flex;flex-direction:column;gap:6px">
            <div class="skeleton" style="height:11px;width:75%"></div>
            <div class="skeleton" style="height:9px;width:35%"></div>
          </div>
        </div>
      </template>

      <div v-else-if="!notifications.length" class="card empty-state">
        <div class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M13.73 21a2 2 0 0 1-3.46 0m10.19-2.28L4 3m14 5a6 6 0 0 0-9.33-5m2.09 13.73A6 6 0 0 1 6 8v5.58L3 19h15.28"></path></svg>
        </div>
        <h3>All caught up</h3>
        <p>No notifications yet</p>
      </div>

      <template v-else>
        <div class="section-label" v-if="unread.length">New · {{ unread.length }}</div>
        <div :class="['card notif-row', 'unread']" v-for="n in unread" :key="n.id">
          <div class="notif-icon">
            <svg v-if="n.type === 'follow_request'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <svg v-else-if="n.type === 'new_follower' || n.type === 'follow_accept'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            <svg v-else-if="n.type === 'group_invite' || n.type === 'group_added' || n.type === 'group_accept'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            <svg v-else-if="n.type === 'group_join_request'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <svg v-else-if="n.type === 'group_event_created'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            <svg v-else-if="n.type === 'like'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
            <svg v-else-if="n.type === 'comment'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
            <svg v-else-if="n.type === 'message'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9zm-8 13a2 2 0 0 0 4 0"></path></svg>
          </div>
          <div style="flex:1">
            <p style="font-size:0.9rem;color:var(--text);margin-bottom:6px;cursor:pointer" @click="navigateTo(n)">{{ n.message }}</p>
            <span style="font-size:0.74rem;color:var(--text-3)">{{ relTime(n.created_at) }}</span>

            <!-- Follow request actions -->
            <div v-if="n.type === 'follow_request' && !n.actioned" style="display:flex;gap:8px;margin-top:10px" @click.stop>
              <button class="btn btn-primary btn-sm" :disabled="n.acting" @click="acceptFollow(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'accept'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Accept
              </button>
              <button class="btn btn-ghost btn-sm" :disabled="n.acting" @click="declineFollow(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'decline'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                Decline
              </button>
            </div>
            <div v-if="n.type === 'follow_request' && n.actioned" style="margin-top:8px;font-size:0.8rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
              <svg v-if="n.actioned === 'accepted'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              {{ n.actioned === 'accepted' ? 'Accepted' : 'Declined' }}
            </div>

            <!-- Group join request actions (shown to group creator) -->
            <div v-if="n.type === 'group_join_request' && !n.actioned" style="display:flex;gap:8px;margin-top:10px" @click.stop>
              <button class="btn btn-primary btn-sm" :disabled="n.acting" @click="acceptGroupJoinRequest(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'accept'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Accept
              </button>
              <button class="btn btn-ghost btn-sm" :disabled="n.acting" @click="declineGroupJoinRequest(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'decline'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                Decline
              </button>
            </div>
            <div v-if="n.type === 'group_join_request' && n.actioned" style="margin-top:8px;font-size:0.8rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
              <svg v-if="n.actioned === 'accepted'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              {{ n.actioned === 'accepted' ? 'Accepted' : 'Declined' }}
            </div>

            <!-- Group invite actions -->
            <div v-if="n.type === 'group_invite' && !n.actioned" style="display:flex;gap:8px;margin-top:10px" @click.stop>
              <button class="btn btn-primary btn-sm" :disabled="n.acting" @click="acceptGroupInvite(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'accept'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Accept
              </button>
              <button class="btn btn-ghost btn-sm" :disabled="n.acting" @click="declineGroupInvite(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'decline'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                Decline
              </button>
            </div>
            <div v-if="n.type === 'group_invite' && n.actioned" style="margin-top:8px;font-size:0.8rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
              <svg v-if="n.actioned === 'accepted'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              {{ n.actioned === 'accepted' ? 'Joined group' : 'Declined' }}
            </div>
          </div>
          <div v-if="!n.is_read" class="notif-dot"></div>
        </div>

        <div class="section-label" style="margin-top:20px" v-if="read.length">Earlier</div>
        <div class="card notif-row" v-for="n in read" :key="n.id" style="cursor:pointer" @click="navigateTo(n)">
          <div class="notif-icon" style="opacity:0.5">
            <svg v-if="n.type === 'follow_request'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <svg v-else-if="n.type === 'new_follower' || n.type === 'follow_accept'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            <svg v-else-if="n.type === 'group_invite' || n.type === 'group_added' || n.type === 'group_accept'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            <svg v-else-if="n.type === 'group_join_request'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <svg v-else-if="n.type === 'group_event_created'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            <svg v-else-if="n.type === 'like'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
            <svg v-else-if="n.type === 'comment'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
            <svg v-else-if="n.type === 'message'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9zm-8 13a2 2 0 0 0 4 0"></path></svg>
          </div>
          <div style="flex:1">
            <p style="font-size:0.9rem;color:var(--text-2);margin-bottom:3px">{{ n.message }}</p>
            <span style="font-size:0.74rem;color:var(--text-3)">{{ relTime(n.created_at) }}</span>

            <!-- Still show action buttons for pending requests in the read section -->
            <div v-if="n.type === 'follow_request' && !n.actioned" style="display:flex;gap:8px;margin-top:10px" @click.stop>
              <button class="btn btn-primary btn-sm" :disabled="n.acting" @click="acceptFollow(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'accept'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Accept
              </button>
              <button class="btn btn-ghost btn-sm" :disabled="n.acting" @click="declineFollow(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'decline'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                Decline
              </button>
            </div>
            <div v-if="n.type === 'follow_request' && n.actioned" style="margin-top:8px;font-size:0.8rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
              <svg v-if="n.actioned === 'accepted'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              {{ n.actioned === 'accepted' ? 'Accepted' : 'Declined' }}
            </div>

            <div v-if="n.type === 'group_join_request' && !n.actioned" style="display:flex;gap:8px;margin-top:10px" @click.stop>
              <button class="btn btn-primary btn-sm" :disabled="n.acting" @click="acceptGroupJoinRequest(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'accept'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Accept
              </button>
              <button class="btn btn-ghost btn-sm" :disabled="n.acting" @click="declineGroupJoinRequest(n)" style="display:inline-flex;align-items:center;gap:4px">
                <svg v-if="n.acting !== 'decline'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                Decline
              </button>
            </div>
            <div v-if="n.type === 'group_join_request' && n.actioned" style="margin-top:8px;font-size:0.8rem;color:var(--text-3);display:inline-flex;align-items:center;gap:4px">
              <svg v-if="n.actioned === 'accepted'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              {{ n.actioned === 'accepted' ? 'Accepted' : 'Declined' }}
            </div>
          </div>
        </div>
      </template>
    </div>

  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './Navbar.vue'
import { apiFetch } from '../services/api'
import { onNotifEvent, resetNotifCount } from '../stores/notif'

const router = useRouter()
const notifications = ref([])
const loading = ref(true)
const markingAll = ref(false)
let unsubNotif = null

const unread = computed(() => notifications.value.filter(n => !n.is_read))
const read = computed(() => notifications.value.filter(n => n.is_read))



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

function groupJoinParts(n) {
  // New format: "requesterUserId:groupId". Old DB records may have just groupId.
  const ref = n.reference_id || ''
  const colon = ref.indexOf(':')
  if (colon === -1) return { userId: '', groupId: ref }
  return { userId: ref.slice(0, colon), groupId: ref.slice(colon + 1) }
}

function navigateTo(n) {
  if (n.type === 'follow_request' || n.type === 'new_follower' || n.type === 'follow_accept') {
    if (n.reference_id) router.push('/profile/' + n.reference_id)
  } else if (n.type === 'group_join_request') {
    const { groupId } = groupJoinParts(n)
    if (groupId) router.push('/groups/' + groupId)
  } else if (n.type === 'group_invite' || n.type === 'group_added' || n.type === 'group_accept' || n.type === 'group_event_created') {
    if (n.reference_id) router.push('/groups/' + n.reference_id)
  } else if (n.type === 'comment' || n.type === 'like') {
    router.push('/feed')
  }
  if (!n.is_read) markOne(n)
}

async function load() {
  try {
    const d = await apiFetch('/notifications')
    const prev = Object.fromEntries((notifications.value).map(n => [n.id, { actioned: n.actioned, acting: n.acting }]))
    notifications.value = (Array.isArray(d) ? d : []).map(n => ({ ...n, ...(prev[n.id] || {}) }))
  } catch {} finally { loading.value = false }
}

async function markOne(n) {
  if (n.is_read) return
  try { await apiFetch('/notifications/read', { method:'POST', body:JSON.stringify({ notification_id:n.id }) }); n.is_read = true } catch {}
}

async function markAll() {
  markingAll.value = true
  try {
    await apiFetch('/notifications/read-all', { method:'POST' })
    notifications.value.forEach(n => n.is_read = true)
    resetNotifCount()
  } catch {} finally { markingAll.value = false }
}

async function clearAll() {
  if (!confirm('Delete all notifications? This cannot be undone.')) return
  try {
    await apiFetch('/notifications/clear-all', { method:'DELETE' })
    notifications.value = []
    resetNotifCount()
  } catch {}
}

async function acceptGroupJoinRequest(n) {
  n.acting = 'accept'
  try {
    const { userId, groupId } = groupJoinParts(n)
    await apiFetch('/groups/request/accept', { method:'POST', body:JSON.stringify({ group_id: groupId, user_id: userId }) })
    n.actioned = 'accepted'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

async function declineGroupJoinRequest(n) {
  n.acting = 'decline'
  try {
    const { userId, groupId } = groupJoinParts(n)
    await apiFetch('/groups/request/decline', { method:'POST', body:JSON.stringify({ group_id: groupId, user_id: userId }) })
    n.actioned = 'declined'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

async function acceptFollow(n) {
  n.acting = 'accept'
  try {
    await apiFetch('/follow/accept', { method:'POST', body:JSON.stringify({ request_id: n.reference_id }) })
    n.actioned = 'accepted'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

async function declineFollow(n) {
  n.acting = 'decline'
  try {
    await apiFetch('/follow/decline', { method:'POST', body:JSON.stringify({ request_id: n.reference_id }) })
    n.actioned = 'declined'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

async function acceptGroupInvite(n) {
  n.acting = 'accept'
  try {
    await apiFetch('/groups/invite/accept', { method:'POST', body:JSON.stringify({ group_id: n.reference_id }) })
    n.actioned = 'accepted'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

async function declineGroupInvite(n) {
  n.acting = 'decline'
  try {
    await apiFetch('/groups/invite/decline', { method:'POST', body:JSON.stringify({ group_id: n.reference_id }) })
    n.actioned = 'declined'; n.is_read = true
  } catch (e) { alert(e.message) } finally { n.acting = null }
}

onMounted(async () => {
  try { await apiFetch('/me') } catch { router.push('/login'); return }
  await load()
  // Mark all as read when entering the page — but skip actionable requests
  // so they stay in the "New" section until the user acts on them
  try {
    await apiFetch('/notifications/read-all', { method:'POST' })
    notifications.value.forEach(n => {
      if (n.type !== 'follow_request' && n.type !== 'group_join_request') {
        n.is_read = true
      }
    })
    resetNotifCount()
  } catch {}
  // Subscribe to live WS push events
  unsubNotif = onNotifEvent((event) => {
    // Avoid duplicates if the WS delivers an event for a notification already loaded from REST
    if (event.id && notifications.value.some(n => n.id === event.id)) return
    notifications.value.unshift({
      id: event.id,
      type: event.type,
      reference_id: event.reference_id,
      message: event.message,
      is_read: false,
      created_at: new Date().toISOString(),
    })
  })
})

onUnmounted(() => {
  if (unsubNotif) unsubNotif()
})
</script>

<style scoped>
.notif-row {
  display:flex; align-items:flex-start; gap:14px;
  padding:14px 18px; margin-bottom:8px;
  border-radius:var(--radius); border:1px solid var(--border);
  transition:border-color 0.15s;
}
.notif-row.unread { border-color:rgba(0,122,255,0.2); background:rgba(0,122,255,0.04); }
.notif-row.unread:hover { border-color:var(--accent); }
.notif-icon { width:36px; height:36px; display:flex; align-items:center; justify-content:center; background:var(--surface2); border-radius:50%; flex-shrink:0; color:var(--text); }
.notif-dot { width:8px; height:8px; border-radius:50%; background:var(--accent); flex-shrink:0; margin-top:6px; }

</style>
