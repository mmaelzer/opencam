import {
  contains,
  first,
  lower
} from './utils'

export function elements (str) {
  let div = document.createElement('div')
  div.innerHTML = str.trim()
  return div.childNodes
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