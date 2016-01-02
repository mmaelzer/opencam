export function mapStep (arr, count, fn) {
  let results = []
  for (let i = 0; i < arr.length; i += count) {
    results = results.concat(fn(arr.slice(i, i + count)))
  }
  return results
}

export function each (arr, fn) {
  return Array.prototype.forEach.call(arr, fn)
}

export function first (arr) {
  return arr[0]
}

export function last (arr) {
  return arr[arr.length - 1]
}

export function split (str, del) {
  return String.prototype.split.call(str, del)
}

export function map (arr, fn) {
  return Array.prototype.map.call(arr, fn)
}

export function join (arr, str) {
  str = arguments.length === 1 ? '' : str
  return Array.prototype.join.call(arr, str)
}

export function reduce (arr, fn, base) {
  return Array.prototype.reduce.call(arr, fn, base)
}

export function flatten (arr) {
  return reduce(arr, i => {
    return Array.isArray(i) ? i.concat(i) : i
  }, [])
}
