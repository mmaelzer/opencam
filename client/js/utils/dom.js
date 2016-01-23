export function element (str) {
  let div = document.createElement('div')
  div.innerHTML = str.trim()
  return div.childNodes
}

export function on (el, event, listener, useCapture = false) {
  return el.addEventListener(event, listener, useCapture)
}

export function off (el, event, listener) {
  return el.removeEventListener(event, listener)
}

export function prepend (parent, child) {
  if (parent.firstChild) {
    parent.insertBefore(child, parent.firstChild)
  } else {
    parent.appendChild(child)
  }
  return child
}

export function remove (el) {
  return el.parentElement.removeChild(el)
}

export function findClass (el, className) {
  if (!el) return null
  if (el.className === className) return el
  for (let i = 0; i < el.childNodes.length; i++) {
    if (el.childNodes[i].className === className) {
      return el.childNodes[i]
    } else {
      let found = findClass(el.childNodes[i], className)
      if (found) return found
    }
  }
  return null
}
