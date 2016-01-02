import superagent from 'superagent'
import superagentJson from 'superagent-jsonapify'
import eventTemplate from '../templates/event'
import rowTemplate from '../templates/row'
import {
  after,
  captureErr,
  each,
  first,
  join,
  last,
  map,
  partial,
  split
} from '../utils/utils'
superagentJson(superagent)

const request = superagent

function getTimeout (image1, image2) {
  if (typeof image1 === 'undefined' || typeof image2 === 'undefined') {
    return 0
  }
  let time1 = first(split(last(split(image1.src, '/')), '.'))
  let time2 = first(split(last(split(image2.src, '/')), '.'))
  return (Number(time2) - Number(time1)) / 1000 / 1000
}

function getEventWithFrames (id, callback) {
  request.get('/api/event/' + id).end((err, res) => {
    if (err) return callback(err)
    callback(null, res.body)
  })
}

function getEvents (page, callback) {
  const url = '/api/events?page=' + page
  request.get(url).end(captureErr((res) => {
    callback(null, Array.isArray(res.body) ? res.body : [])
  }))
}

function buildImageForFrame (frame) {
  const image = new window.Image()
  image.src = frame
  return image
}

function play (frames, div) {
  return function () {
    while (div.hasChildNodes()) {
      div.removeChild(div.lastChild)
    }
    function loadFrame (i) {
      if (i === frames.length) return
      if (i > 0) div.removeChild(frames[i - 1])
      div.appendChild(frames[i])
      const timeout = getTimeout(frames[i], frames[i + 1])
      console.log(timeout)
      window.setTimeout(partial(loadFrame, i + 1), timeout)
    }
    loadFrame(0)
  }
}

function playEvent (el) {
  return function (event) {
    let imgContainer = first(el.getElementsByClassName('event-image-container'))
    let frames = map(event.frames, buildImageForFrame)
    let playWhenReady = after(frames.length, play(frames, imgContainer))
    each(frames, frame => frame.onload = playWhenReady)
  }
}

function onEventClick (e) {
  e.preventDefault()
  const el = e.currentTarget
  getEventWithFrames(el.dataset.id, captureErr(playEvent(el)))
}

function buildEventEl (el) {
  el.addEventListener('click', onEventClick, true)
  return () => el.removeEventListener('click', onEventClick)
}

function renderEvents (events) {
  const eventsDiv = document.getElementById('events')
  eventsDiv.innerHTML = rowTemplate(join(map(events, eventTemplate)))
  const eventsEl = document.getElementsByClassName('event')
  return map(eventsEl, buildEventEl)
}

(function main () {
  getEvents(1, (err, events) => {
    if (err) return console.error(err)
    renderEvents(events)
  })
})()
