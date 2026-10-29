import { test } from 'node:test';
import assert from 'node:assert/strict';
import { createElement as h, createStore, createRouter, render } from '../framework/index.js';

test('element descriptions flatten nested children and preserve text and zero', () => {
  assert.deepEqual(h('div', { id: 'test' }, ['text', [0, false, true]]), {
    tag: 'div', attrs: { id: 'test' }, children: ['text', 0],
  });
});

test('shared store updates, suppresses identical values and unsubscribes', () => {
  const store = createStore({ count: 0 });
  const observed = [];
  const stop = store.subscribe(state => observed.push(state.count));
  store.setState(state => ({ count: state.count + 1 }));
  store.setState(store.getState());
  stop();
  store.setState({ count: 2 });
  assert.deepEqual(observed, [1]);
  assert.equal(store.getState().count, 2);
  assert.throws(() => store.subscribe(null), TypeError);
});

test('router responds to navigation and browser history events and cleans up', () => {
  const events = new EventTarget();
  globalThis.window = {
    location: { hash: '' },
    addEventListener: events.addEventListener.bind(events),
    removeEventListener: events.removeEventListener.bind(events),
  };
  const router = createRouter({ '/': 'all', '/active': 'active' });
  const seen = [];
  router.subscribe(route => seen.push(route));
  router.navigate('active');
  assert.equal(window.location.hash, '/active');
  window.location.hash = '#/active';
  events.dispatchEvent(new Event('hashchange'));
  assert.deepEqual(seen, [{ path: '/active', route: 'active' }]);
  window.location.hash = '#/unknown';
  assert.equal(router.resolve().route, 'all');
  router.destroy();
  events.dispatchEvent(new Event('hashchange'));
  assert.equal(seen.length, 1);
  delete globalThis.window;
});

test('renderer creates nested DOM, properties, attributes and working declarative events', () => {
  class Element {
    constructor(tag) { this.tag = tag; this.attrs = {}; this.children = []; this.events = {}; }
    setAttribute(key, value) { this.attrs[key] = String(value); }
    appendChild(child) { this.children.push(child); }
    addEventListener(name, callback) { this.events[name] = callback; }
    set innerHTML(value) { this.children = []; }
  }
  globalThis.Element = Element;
  globalThis.document = { createElement: tag => new Element(tag), createTextNode: text => ({ text }) };
  const root = new Element('main');
  let clicks = 0;
  render(h('div', { className: 'card', 'aria-hidden': false },
    h('button', { onClick: () => clicks++, disabled: false }, 'Save'),
    h('input', { checked: true, value: 'hello' })), root);
  const card = root.children[0];
  assert.equal(card.attrs.class, 'card');
  assert.equal(card.attrs['aria-hidden'], 'false');
  card.children[0].events.click();
  assert.equal(clicks, 1);
  assert.equal(card.children[0].attrs.disabled, undefined);
  assert.equal(card.children[1].checked, true);
  assert.equal(card.children[1].value, 'hello');
  render(h('p', {}, 'Replacement'), root);
  assert.equal(root.children.length, 1);
  assert.equal(root.children[0].tag, 'p');
  assert.throws(() => render(h('p'), null), TypeError);
  delete globalThis.Element;
  delete globalThis.document;
});
