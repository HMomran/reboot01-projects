import { setLoggedIn } from '../stores/auth'

export const BACKEND_URL = ''

export async function apiFetch(path, options = {}) {
    const isFormData = options.body instanceof FormData

    const response = await fetch(`/api${path}`, {
        credentials: 'include',
        ...options,
        headers: {
            ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
            ...(options.headers || {}),
        },
    })

    const contentType = response.headers.get('content-type') || ''
    const data = contentType.includes('application/json')
        ? await response.json()
        : await response.text()

    if (response.status === 401) {
        const publicPaths = ['/login', '/register', '/feed']
        const isPasswordChange = path.includes('/profile/password')
        if (!publicPaths.some(p => window.location.pathname.startsWith(p)) && !isPasswordChange) {
            setLoggedIn(false)
            sessionStorage.setItem('session_msg', 'Your session has expired. Please log in again.')
            window.location.href = '/login'
        }
        throw new Error(typeof data === 'string' ? data : data.error || 'Unauthorized')
    }

    if (!response.ok) {
        throw new Error(typeof data === 'string' ? data : data.error || 'Request failed')
    }

    return data
}
