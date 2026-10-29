export function createStore(initialState) {
  let state = initialState;
  const listeners = new Set();

  function getState() {
    return state;
  }

  function setState(update) {
    const nextState =
      typeof update === "function" ? update(state) : update;

    if (Object.is(nextState, state)) return state;

    state = nextState;
    for (const listener of [...listeners]) listener(state);
    return state;
  }

  function subscribe(listener) {
    if (typeof listener !== "function") {
      throw new TypeError("subscribe expects a function");
    }

    listeners.add(listener);
    return () => listeners.delete(listener);
  }

  return { getState, setState, subscribe };
}
