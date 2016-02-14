import superagent from 'superagent'
import superagentJson from 'superagent-jsonapify'
import eventTemplate from '../templates/event'
import rowTemplate from '../templates/row'
import {
  hasClass,
  off,
  on,
  toggleClass
} from '../utils/dom'
import {
  _,
  after,
  captureErr,
  del,
  each,
  filter,
  first,
  join,
  last,
  map,
  partial,
  pluck,
  split
} from '../utils/utils'
superagentJson(superagent)

const request = superagent
const players = {}
const images = {}

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

function getEventsForCameras (cameras, page, callback) {
  let cameraQuery = join(map(cameras, (c) => 'camera_id=' + c), '&')
  const url = '/api/events?page=' + page + '&' + cameraQuery
  request.get(url).end(captureErr((res) => {
    callback(null, Array.isArray(res.body) ? res.body : [])
  }))
}

function buildImageForFrame (frame) {
  const image = new window.Image()
  image.src = frame
  image.className = 'event-image'
  return image
}

function play (eventId, frames, container) {
  const info = first(
    container.getElementsByClassName('event-info')
  )
  const progress = first(
    container.getElementsByClassName('event-video-progress')
  )
  progress.style.width = '0%'
  info.style.opacity = 0
  function loadFrame (i) {
    if (i === frames.length) {
      info.style.opacity = 1
      del(players, eventId)
      return
    }
    let img = first(container.getElementsByClassName('event-image'))
    container.replaceChild(frames[i], img)
    progress.style.width = ((100 * (i + 1)) / frames.length) + '%'
    const timeout = getTimeout(frames[i], frames[i + 1])
    players[eventId] = window.setTimeout(partial(loadFrame, i + 1), timeout)
  }
  loadFrame(0)
}

function playEvent (el) {
  return function (event) {
    let container = first(el.getElementsByClassName('event-image')).parentElement
    // If playing, stop the player
    if (event.id in players) {
      window.clearTimeout(players[event.id])
      del(players, event.id)
      return
    }
    // If images cached, used cached images instead of
    // fetching new images from the server
    if (event.id in images) {
      play(event.id, images[event.id], container)
      return
    }
    let frames = map(event.frames, buildImageForFrame)
    let loaded = after(frames.length, () => {
      images[event.id] = frames
      play(event.id, frames, container)
    })
    each(frames, frame => frame.onload = loaded)
  }
}

function onEventClick (e) {
  e.preventDefault()
  const el = e.currentTarget
  getEventWithFrames(el.dataset.id, captureErr(playEvent(el)))
}

function buildEventEl (el) {
  on(el, 'click', onEventClick)
  return partial(off, el, 'click', onEventClick)
}

function renderEvents (events) {
  const cameras = window.OPENCAM_CAMERAS
  const eventsDiv = document.getElementById('events')
  eventsDiv.innerHTML = rowTemplate(
    join(map(events, partial(eventTemplate, _, cameras)))
  )
  const eventsEl = document.getElementsByClassName('event')
  return map(eventsEl, buildEventEl)
}

function bindFilters () {
  let filters = document.getElementsByClassName('filter-camera')
  map(filters, (cam) => {
    on(cam, 'click', (e) => {
      e.preventDefault()
      toggleClass(cam, 'active')
      let activeCameraIds = map(
        filter(
          filters,
          partial(hasClass, _, 'active')
        ),
        (c) => c.dataset.id
      )
      getEventsForCameras(activeCameraIds, 1, captureErr(renderEvents))
    })
  })
}

(function main () {
  getEvents(1, captureErr(renderEvents))
  bindFilters()
})()
