import moment from 'moment'
import superagent from 'superagent'
import superagentJson from 'superagent-jsonapify'
import eventTemplate from '../templates/event'
import rowTemplate from '../templates/row'
import {
  addClass,
  hasClass,
  off,
  on,
  removeClass,
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
  flatten,
  group,
  invoke,
  join,
  last,
  map,
  partial,
  pluck,
  parallel,
  reduce,
  sort,
  split,
  values
} from '../utils/utils'
superagentJson(superagent)

const request = superagent
const players = {}
const images = {}
let EVENTS = []
let unbindListeners = []

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
    callback(err, res.body)
  })
}

function getEvents (page, callback) {
  fetchEvents(`/api/events?page=${page}`, callback)
}

function fetchEvents (url, callback) {
  request.get(url).end((err, res) => {
    EVENTS = Array.isArray(res.body) ? res.body : []
    callback(err, EVENTS)
  })
}

function getEventsForCameras (cameras, page, callback) {
  let cameraQuery = join(map(cameras, (c) => 'camera_id=' + c), '&')
  fetchEvents(`/api/events?page=${page}&${cameraQuery}`, callback)
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
  let id = split(el.dataset.id, ',')
  if (id.length > 1) {
    parallel(id, getEventWithFrames, captureErr((events) => {
      playEvent(el)({
        id: join(id, ','),
        frames: sort(flatten(pluck(events, 'frames')))
      })
    }))
  } else {
    getEventWithFrames(first(id), captureErr(playEvent(el)))
  }
}

function buildEventEl (el) {
  on(el, 'click', onEventClick)
  return partial(off, el, 'click', onEventClick)
}

function renderEvents (events) {
  invoke(unbindListeners)
  const cameras = window.OPENCAM_CAMERAS
  const eventsDiv = document.getElementById('events')
  eventsDiv.innerHTML = rowTemplate(
    join(map(events, partial(eventTemplate, _, cameras)))
  )
  const eventsEl = document.getElementsByClassName('event')
  unbindListeners = map(eventsEl, buildEventEl)
}

function bindCameraFilters () {
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

function bindGroupFilters () {
  let filters = document.getElementsByClassName('filter-group')
  map(filters, g => {
    on(g, 'click', e => {
      e.preventDefault()
      if (hasClass(g, 'active')) return
      map(filters, f => removeClass(f, 'active'))
      addClass(g, 'active')
      let range = +g.dataset.range
      if (range === 0) return renderEvents(EVENTS)
      let grouped = group(EVENTS, (e) => {
        let start = moment(e.start_time).unix()
        let timeGroup = Math.round(start / (range * 60))
        return `${e.camera_id}:${timeGroup}`
      })
      renderEvents(buildGroupEvents(grouped))
    })
  })
}

function buildGroupEvents (grouped) {
  return reduce(values(grouped), (events, grouping) => {
    let toUnix = (time) => {
      return moment(time).valueOf()
    }
    let start = Math.min(...map(pluck(grouping, 'start_time'), toUnix))
    let end = Math.max(...map(pluck(grouping, 'end_time'), toUnix))
    let firstEvent = first(grouping)
    return events.concat({
      id: pluck(grouping, 'id'),
      camera_id: firstEvent.camera_id,
      first_frame: firstEvent.first_frame,
      duration: (end - start) / 1000,
      start_time: moment(start).toISOString(),
      end_time: moment(end).toISOString()
    })
  }, [])
}

function bindFilters () {
  bindCameraFilters()
  bindGroupFilters()
}

(function main () {
  getEvents(1, captureErr(renderEvents))
  bindFilters()
})()
