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

export function breaker (test) {
  if (test()) throw new Breaker()
}

export function contains (iterable, val) {
  return iterable && iterable.indexOf(val) >= 0
}

export function del (obj, prop) {
  delete obj[prop]
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
  return reduce(arr, (flat, i) => flat.concat(array(i)), [])
}

export function invoke (arr, var_args) {
  return map(arr, fn => fn.call(this, ...slice(arguments, 1)))
}

export function isFunction (o) {
  return typeof o === 'function'
}

export function group (arr, prop) {
  return reduce(arr, (obj, item) => {
    let key = result(item, prop)
    if (key in obj) {
      obj[key].push(item)
    } else {
      obj[key] = [item]
    }
    return obj
  }, {})
}

export function join (arr, str = '') {
  return Array.prototype.join.call(arr, str)
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

export function parallel (arr, iterator, callback = _) {
  let results = []
  let called = 0
  for (let i = 0; i < arr.length; i++) iterate(i)
  function iterate (index) {
    iterator(arr[index], (err, result) => {
      if (called === -1) return
      if (err) {
        called = -1
        return callback(err)
      }
      results[index] = result
      if (++called === arr.length) callback(undefined, results)
    }, index)
  }
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
      } catch (e) {
        if (e instanceof Breaker) return
        throw e
      }
    }
    return res
  }
}

export function pluck (arr, prop) {
  return map(arr, partial(result, _, prop))
}

export function reduce (arr, fn, base) {
  return Array.prototype.reduce.call(arr, fn, base)
}

export function result (obj, getter) {
  return isFunction(getter) ? getter(obj) : obj[getter]
}

export function propSorter (prop) {
  return (a, b) => {
    let ra = result(a, prop)
    let rb = result(b, prop)
    if (ra > rb) return -1
    if (ra < rb) return 1
    return 0
  }
}

export function sort (arr, prop) {
  let fn = typeof prop === 'string' ? propSorter(prop) : prop
  return Array.prototype.sort.call(arr, fn)
}

export function values (dict) {
  return map(Object.keys(dict), key => dict[key])
}

export function without (arr, val) {
  let test = typeof val === 'function' ? val : (v) => v !== val
  return filter(arr, test)
}

class Breaker extends Error {}
