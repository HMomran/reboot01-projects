export function createRouter(routes, fallback = "/") {
  const listeners = new Set();

  function getPath() {
    return window.location.hash.slice(1) || fallback;
  }

  function resolve() {
    const path = getPath();
    return {
      path,
      route: routes[path] ?? routes[fallback] ?? null,
    };
  }

  function notify() {
    const current = resolve();
    for (const listener of [...listeners]) listener(current);
  }

  function navigate(path) {
    const normalized = path.startsWith("/") ? path : `/${path}`;
    if (getPath() === normalized) notify();
    else window.location.hash = normalized;
  }

  function subscribe(listener) {
    listeners.add(listener);
    return () => listeners.delete(listener);
  }

  window.addEventListener("hashchange", notify);

  function destroy() {
    window.removeEventListener("hashchange", notify);
    listeners.clear();
  }

  return { getPath, resolve, navigate, subscribe, destroy };
}
