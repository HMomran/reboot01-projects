# Mini Framework

A small, dependency-free JavaScript framework with a virtual element tree, declarative events, shared state, and hash routing. The included TodoMVC demonstrates all four features.

## Run it

```bash
npm start
```

Open the printed URL. The project runs from its root and requires no dependencies or build step.

## API

Import the public functions from `framework/index.js`:

```js
import { createElement, render, createStore, createRouter } from "./framework/index.js";
```

### Create and nest elements

`createElement(tag, attributes, ...children)` returns a plain JavaScript description of a DOM node. Children can be more element descriptions, strings, numbers, or nested arrays.

```js
const card = createElement(
  "div",
  { class: "card", id: "welcome" },
  createElement("h2", {}, "Hello"),
  createElement("p", {}, "This paragraph is nested inside the card."),
);

render(card, document.getElementById("app"));
```

Attributes are expressed as an object. Use `class` or `className` for classes and `htmlFor` for a label's `for` attribute. Boolean attributes such as `disabled` and `checked` accept booleans.

### Handle events

Pass a function under an `on<EventName>` attribute:

```js
const button = createElement(
  "button",
  { onClick: () => console.log("clicked"), disabled: false },
  "Save",
);
```

The framework translates this declarative event property into the browser event listener. Application code describes *what* should happen alongside the element instead of calling `addEventListener` itself. Any browser event works, such as `onInput`, `onKeyDown`, `onSubmit`, or `onScroll`.

### Manage state

```js
const counter = createStore({ count: 0 });

counter.subscribe((state) => {
  render(createElement("p", {}, `Count: ${state.count}`), app);
});

counter.setState((state) => ({ ...state, count: state.count + 1 }));
console.log(counter.getState());
```

`getState()` reads the current shared value. `setState(valueOrUpdater)` replaces it and notifies subscribers. `subscribe(listener)` returns an unsubscribe function. State lives in the store rather than a particular page, so every route or component can use the same source of truth.

### Route by URL

```js
const router = createRouter({
  "/": "home",
  "/about": "about",
});

router.subscribe(({ path, route }) => {
  console.log(path, route);
});

router.navigate("/about"); // URL becomes #/about
```

The router uses the URL hash, so it works with a static file server and browser back/forward navigation. `getPath()` reads the current path, `resolve()` returns its route value, and `destroy()` removes the router's browser listener.

## Why it works this way

Element descriptions are plain objects, which separates the desired UI from direct DOM manipulation. `render` recursively turns that tree into real nodes and attaches attributes and events. A full render is intentionally used here: it keeps this learning framework predictable and small; a production framework could replace it with tree diffing without changing how application views are written.

The store owns state and calls subscribers after an immutable update. The TodoMVC subscriber persists todos and renders a new view, demonstrating inversion of control: application code declares a view and update functions, while the framework decides when to notify and create DOM nodes. The router similarly converts browser URL changes into application notifications.

## TodoMVC behavior

The example supports adding, completing, editing (double-click), deleting, toggling all, clearing completed items, persistence through `localStorage`, and All/Active/Completed URL filters. Its implementation is in `todo/app.js` and uses only this framework and browser APIs.

## Structure

```text
framework/
  createElement.js  virtual element descriptions
  render.js         DOM creation and declarative events
  store.js          shared observable state
  router.js         hash URL synchronization
  index.js          public exports
todo/
  app.js            TodoMVC application
  style.css         TodoMVC presentation
index.html          root entry point
server.js           zero-dependency development server
```

## Verification

Run `npm test` (Node.js 22 or newer). These dependency-free checks cover element
nesting, store notifications, URL routing, DOM creation and declarative events.
The renderer test uses a small DOM substitute; browser interaction still needs
manual verification.

For a browser acceptance check:

1. Add two todos with Enter, including a title with surrounding spaces. Verify
   trimming, the item count, and that you can immediately type another todo.
2. Toggle one item, visit Active and Completed, then use browser Back/Forward.
3. Double-click a title, save with Enter, cancel with Escape, and save on blur.
   Saving an empty title should delete that todo.
4. Toggle all, clear completed, and delete individual items.
5. Reload on `#/active` and verify both saved todos and the selected filter.

Rendering replaces the container's children, so applications that need to keep
input drafts or focus across updates must restore them, as `draw` does in the
TodoMVC example. Store and router subscriptions notify on changes; invoke your
initial view explicitly. Treat store values as immutable and return new objects
from updater functions. Browser storage failures are tolerated, but persistence
requires available local storage.

## Project documentation

[Project guide: purpose, architecture, setup, and usage](PROJECT_GUIDE.md)
