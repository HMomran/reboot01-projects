function createDOMNode(virtualNode) {
  if (
    typeof virtualNode === "string" ||
    typeof virtualNode === "number"
  ) {
    return document.createTextNode(String(virtualNode));
  }

  if (virtualNode === null || virtualNode === undefined) {
    return document.createTextNode("");
  }

  const element = document.createElement(virtualNode.tag);

  setAttributes(element, virtualNode.attrs);

  for (const child of virtualNode.children || []) {
    const childNode = createDOMNode(child);
    element.appendChild(childNode);
  }

  return element;
}

function setAttributes(element, attrs) {
  for (const [name, value] of Object.entries(attrs || {})) {
    if (name.startsWith("on") && typeof value === "function") {
      const eventName = name.slice(2).toLowerCase();

      element.addEventListener(eventName, value);
      continue;
    }

    if (name === "className") {
      element.setAttribute("class", value);
      continue;
    }

    if (name === "htmlFor") {
      element.setAttribute("for", value);
      continue;
    }

    if (name === "value" || name === "checked" || name === "selected") {
      element[name] = value;
      continue;
    }

    if (name.startsWith("aria-") || name.startsWith("data-")) {
      if (value !== null && value !== undefined) element.setAttribute(name, String(value));
      continue;
    }

    if (value === false || value === null || value === undefined) {
      continue;
    }

    if (value === true) {
      element.setAttribute(name, "");
      continue;
    }

    element.setAttribute(name, value);
  }
}

export function render(virtualNode, container) {
  if (!(container instanceof Element)) {
    throw new TypeError("render expects a DOM element as its container");
  }

  container.innerHTML = "";

  const realDOM = createDOMNode(virtualNode);

  container.appendChild(realDOM);
}
