import {
  createElement as h,
  createRouter,
  createStore,
  render,
} from "../framework/index.js";

const STORAGE_KEY = "mini-framework-todos";
const validFilters = new Set(["/", "/active", "/completed"]);

function loadTodos() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY));
    if (!Array.isArray(saved)) return [];
    const ids = new Set();
    return saved.filter((todo) => {
      if (!todo || typeof todo.id !== "string" || !/^[a-zA-Z0-9-]+$/.test(todo.id) ||
          ids.has(todo.id) || typeof todo.title !== "string" || !todo.title.trim() ||
          typeof todo.completed !== "boolean") return false;
      ids.add(todo.id);
      return true;
    });
  } catch {
    return [];
  }
}

const store = createStore({
  todos: loadTodos(),
  filter: "/",
  editingId: null,
});

const router = createRouter(
  { "/": "all", "/active": "active", "/completed": "completed" },
  "/",
);

function updateTodos(change) {
  store.setState((state) => ({
    ...state,
    todos: change(state.todos),
    editingId: null,
  }));
}

function addTodo(event) {
  if (event.key !== "Enter") return;
  const title = event.target.value.trim();
  if (!title) return;

  event.target.value = "";
  updateTodos((todos) => [
    ...todos,
    { id: crypto.randomUUID(), title, completed: false },
  ]);
}

function finishEditing(id, value) {
  if (store.getState().editingId !== id) return;
  const title = value.trim();
  updateTodos((todos) =>
    title
      ? todos.map((todo) => (todo.id === id ? { ...todo, title } : todo))
      : todos.filter((todo) => todo.id !== id),
  );
}

function todoItem(todo, editingId) {
  const classes = [todo.completed && "completed", editingId === todo.id && "editing"]
    .filter(Boolean)
    .join(" ");

  return h(
    "li",
    { class: classes, "data-id": todo.id },
    h(
      "div",
      { class: "view" },
      h("input", {
        class: "toggle",
        type: "checkbox",
        checked: todo.completed,
        "aria-label": `Complete ${todo.title}`,
        onChange: () =>
          updateTodos((todos) =>
            todos.map((item) =>
              item.id === todo.id ? { ...item, completed: !item.completed } : item,
            ),
          ),
      }),
      h(
        "label",
        {
          onDblClick: () => {
            store.setState((state) => ({ ...state, editingId: todo.id }));
            requestAnimationFrame(() => {
              const input = document.querySelector(`li[data-id="${todo.id}"] .edit`);
              input?.focus();
              input?.setSelectionRange(input.value.length, input.value.length);
            });
          },
        },
        todo.title,
      ),
      h("button", {
        class: "destroy",
        "aria-label": `Delete ${todo.title}`,
        onClick: () => updateTodos((todos) => todos.filter((item) => item.id !== todo.id)),
      }),
    ),
    h("input", {
      class: "edit",
      value: todo.title,
      onBlur: (event) => finishEditing(todo.id, event.target.value),
      onKeyDown: (event) => {
        if (event.key === "Enter") event.target.blur();
        if (event.key === "Escape") {
          store.setState((state) => ({ ...state, editingId: null }));
        }
      },
    }),
  );
}

function filterLink(path, label, current) {
  return h(
    "li",
    {},
    h(
      "a",
      {
        href: `#${path}`,
        class: current === path ? "selected" : "",
        onClick: (event) => {
          event.preventDefault();
          router.navigate(path);
        },
      },
      label,
    ),
  );
}

function App(state) {
  const { todos, filter, editingId } = state;
  const remaining = todos.filter((todo) => !todo.completed).length;
  const completed = todos.length - remaining;
  const visible = todos.filter((todo) => {
    if (filter === "/active") return !todo.completed;
    if (filter === "/completed") return todo.completed;
    return true;
  });

  return h(
    "section",
    { class: "todoapp" },
    h(
      "header",
      { class: "header" },
      h("h1", {}, "todos"),
      h("input", {
        class: "new-todo",
        placeholder: "What needs to be done?",
        "aria-label": "New todo",
        autofocus: true,
        onKeyDown: addTodo,
      }),
    ),
    todos.length > 0 &&
      h(
        "main",
        { class: "main" },
        h("input", {
          id: "toggle-all",
          class: "toggle-all",
          type: "checkbox",
          checked: remaining === 0,
          onChange: () => {
            const completed = remaining > 0;
            updateTodos((items) => items.map((todo) => ({ ...todo, completed })));
          },
        }),
        h("label", { htmlFor: "toggle-all" }, "Mark all as complete"),
        h("ul", { class: "todo-list" }, visible.map((todo) => todoItem(todo, editingId))),
      ),
    todos.length > 0 &&
      h(
        "footer",
        { class: "footer" },
        h(
          "span",
          { class: "todo-count" },
          h("strong", {}, remaining),
          ` ${remaining === 1 ? "item" : "items"} left`,
        ),
        h(
          "ul",
          { class: "filters" },
          filterLink("/", "All", filter),
          filterLink("/active", "Active", filter),
          filterLink("/completed", "Completed", filter),
        ),
        completed > 0 &&
          h(
            "button",
            {
              class: "clear-completed",
              onClick: () => updateTodos((items) => items.filter((todo) => !todo.completed)),
            },
            "Clear completed",
          ),
      ),
  );
}

const root = document.getElementById("app");

function draw(state = store.getState()) {
  const input = root.querySelector(".new-todo");
  const wasFocused = document.activeElement === input;
  const draft = input?.value || "";
  const start = input?.selectionStart;
  const end = input?.selectionEnd;
  render(App(state), root);
  const next = root.querySelector(".new-todo");
  next.value = draft;
  if (wasFocused) {
    next.focus();
    next.setSelectionRange(start, end);
  }
}

store.subscribe((state) => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state.todos));
  } catch {
    // The app remains usable when browser storage is unavailable or full.
  }
  draw(state);
});

router.subscribe(({ path }) => {
  store.setState((state) => ({
    ...state,
    filter: validFilters.has(path) ? path : "/",
  }));
});

const initialPath = router.getPath();
if (!validFilters.has(initialPath)) router.navigate("/");
else store.setState((state) => ({ ...state, filter: initialPath }));

draw();
