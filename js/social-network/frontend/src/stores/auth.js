import { ref } from 'vue'

export const isLoggedIn = ref(!!sessionStorage.getItem('logged_in'))

export function setLoggedIn(v) {
  isLoggedIn.value = v
  if (v) sessionStorage.setItem('logged_in', '1')
  else sessionStorage.removeItem('logged_in')
}
