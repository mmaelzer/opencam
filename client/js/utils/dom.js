import {
  contains,
  first,
  join,
  lower,
  split,
  without
} from './utils'

export function elements (str) {
  let div = document.createElement('div')
  div.innerHTML = str.trim()
  return div.childNodes
}

export function empty (el) {
  while (el.hasChildNodes()) {
    el.removeChild(el.lastChild)
  }
  return el
}

export function on (el, event, listener, useCapture = false) {
  el.addEventListener(event, listener, useCapture)
  return el
}

export function off (el, event, listener) {
  el.removeEventListener(event, listener)
  return el
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

export function replace (oldEl, newEl) {
  oldEl.parentElement.replaceChild(newEl, oldEl)
  return newEl
}

export function element (str) {
  return first(elements(str))
}

export function findClass (el, className) {
  if (!el) return null
  if (el.className === className) return el
  for (let i = 0; i < el.childNodes.length; i++) {
    if (contains(el.childNodes[i].className, className)) {
      return el.childNodes[i]
    } else {
      let found = findClass(el.childNodes[i], className)
      if (found) return found
    }
  }
  return null
}

export function findTags (el, tagName) {
  if (!el) return []
  let tags = []
  for (let i = 0; i < el.childNodes.length; i++) {
    if (lower(el.childNodes[i].tagName) === tagName) {
      tags.push(el.childNodes[i])
    } else {
      tags = tags.concat(findTags(el.childNodes[i], tagName))
    }
  }
  return tags
}

export function hasClass (el, className) {
  return contains(el.className, className)
}

export function removeClass (el, className) {
  el.className = join(without(split(el.className, ' '), className), ' ')
  return el
}

export function addClass (el, className) {
  el.className += ' ' + className
  return el
}

export function toggleClass (el, className) {
  return hasClass(el, className) ? removeClass(el, className) : addClass(el, className)
}
