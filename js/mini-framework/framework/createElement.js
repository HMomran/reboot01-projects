export function createElement(tag, attrs = {}, ...children) {
  return {
    tag,
    attrs,
    children: children
      .flat(Infinity)
      .filter((child) => child !== false && child !== true),
  };
}
