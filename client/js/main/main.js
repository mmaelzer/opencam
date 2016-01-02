import superagent from 'superagent'
import superagentJson from 'superagent-jsonapify'
import eventTemplate from '../templates/event'
import rowTemplate from '../templates/row'
import {each, first, join, last, map, split} from '../utils/utils'
superagentJson(superagent)

const request = superagent

function getTimeout (frame1, frame2) {
  let time1 = first(split(last(split(frame1, '/')), '.'))
  let time2 = first(split(last(split(frame2, '/')), '.'))
  return (Number(time2) - Number(time1)) / 1000 / 1000
}

request.get('/api/events').end((err, res) => {
  if (err) return console.error(err)
  const eventsDiv = document.getElementById('events')
  const events = Array.isArray(res.body) ? res.body : []
  eventsDiv.innerHTML = rowTemplate(join(map(events, eventTemplate)))

  const eventsEl = document.getElementsByClassName('event')
  each(eventsEl, el => {
    el.addEventListener('click', e => {
      e.preventDefault()
      request.get('/api/event/' + el.dataset.id).end((err, res) => {
        if (err) return console.error(err)
        let imgContainer = el.getElementsByClassName('event-image-container')[0]
        let img = new window.Image()
        let frames = res.body.frames
        let i = 0
        function loadFrames () {
          if (i === 1) {
            while (imgContainer.hasChildNodes()) {
              imgContainer.removeChild(imgContainer.lastChild)
            }
            imgContainer.appendChild(img)
          }
          if (i >= frames.length) return
          img.src = frames[i++]
          const timeout = i < frames.length - 1
            ? getTimeout(frames[i], frames[i + 1])
            : 0
          img.onload = () => window.setTimeout(loadFrames, timeout)
        }
        loadFrames()
      })
      return false
    }, true)
  })
})
