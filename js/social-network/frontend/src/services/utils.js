const COLORS = ['#f5a623','#3ecf78','#3b8aff','#e04040','#a78bfa','#f472b6']

export function avatarColor(id) {
  let n = 0
  for (let c of (id || '')) n += c.charCodeAt(0)
  return COLORS[n % COLORS.length]
}

export function initials(f, l) {
  return `${(f||'?')[0]}${(l||'?')[0]}`.toUpperCase()
}

export function relTime(s) {
  if (!s) return ''
  const normalized = s.includes('T') ? s : s.replace(' ', 'T') + 'Z'
  const d = Date.now() - new Date(normalized).getTime()
  if (isNaN(d)) return ''
  const sec = Math.floor(d / 1000)
  if (sec < 60) return 'just now'
  const m = Math.floor(sec / 60)
  if (m < 60) return m + 'm ago'
  const h = Math.floor(m / 60)
  if (h < 24) return h + 'h ago'
  return Math.floor(h / 24) + 'd ago'
}
