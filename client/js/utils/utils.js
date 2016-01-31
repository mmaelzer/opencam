export function _ () {}

export function array (arr) {
  return Array.isArray(arr) ? arr : [arr]
}

export function captureErr (fn) {
  return function (err, var_args) {
    if (err) return console.error(err)
    fn.call(this, ...slice(arguments, 1))
  }
}

export function after (times, fn) {
  let calls = 0
  return () => { if (++calls === times) fn() }
}

export function contains (iterable, val) {
  return iterable && iterable.indexOf(val) >= 0
}

export function del (obj, prop) {
  delete obj[result(prop)]
  return obj
}

export function each (arr, fn) {
  return Array.prototype.forEach.call(arr, fn)
}

export function filter (arr, fn) {
  return Array.prototype.filter.call(arr, fn)
}

export function findWhere (arr, predicate) {
  return first(filter(arr, (item) => {
    let keys = Object.keys(predicate)
    for (let i = 0; i < keys.length; i++) {
      let key = keys[i]
      if (item[key] !== predicate[key]) return false
    }
    return true
  }))
}

export function find (arr, fn) {
  return first(filter(arr, fn))
}

export function first (arr) {
  return arr[0]
}

export function flatten (arr) {
  return reduce(arr, i => arr.concat(array(i)), [])
}

export function isFunction (o) {
  return typeof o === 'function'
}

export function group (arr, prop) {
  return reduce(arr, (obj, item) => {
    let key = result(item, prop)
    if (key in obj) {
      if (Array.isArray(obj[key])) {
        obj[key].push(item)
      } else {
        obj[key] = [obj[key], item]
      }
    } else {
      obj[key] = item
    }
    return obj
  }, {})
}

export function last (arr) {
  return arr[arr.length - 1]
}

export function lower (str) {
  return typeof str === 'string' ? str.toLowerCase() : ''
}

export function split (str, del) {
  return String.prototype.split.call(str, del)
}

export function slice (arr, start, end) {
  return Array.prototype.slice.call(arr, start, end)
}

export function map (arr, fn) {
  return Array.prototype.map.call(arr, fn)
}

export function mapStep (arr, count, fn) {
  let results = []
  for (let i = 0; i < arr.length; i += count) {
    results = results.concat(fn(arr.slice(i, i + count)))
  }
  return results
}

export function join (arr, str) {
  str = arguments.length === 1 ? '' : str
  return Array.prototype.join.call(arr, str)
}

export function partial (fn, var_args) {
  const args = slice(arguments, 1)
  return function () {
    let internalArgs = slice(arguments)
    let argsCopy = slice(args)
    for (let i = 0; i < argsCopy.length; i++) {
      if (argsCopy[i] === _) {
        argsCopy[i] = internalArgs.shift()
      }
    }
    return fn.call(this, ...argsCopy.concat(internalArgs))
  }
}

export function pipe (var_args) {
  let args = arguments
  return function () {
    let res
    for (let i = 0; i < args.length; i++) {
      try {
        res = args[i](res)
      } catch (e) { return }
    }
    return res
  }
}

export function breaker (test) {
  if (test()) throw new Error()
}

export function reduce (arr, fn, base) {
  return Array.prototype.reduce.call(arr, fn, base)
}

export function result (obj, getter) {
  return isFunction(getter) ? getter(obj) : obj[getter]
}

