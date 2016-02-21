import moment from 'moment'
import { findWhere } from '../utils/utils'

export default function (event, cameras) {
  const timeString = (time) => moment(time).format('hh:mm:ssa MM/DD/YYYY')
  const time = timeString(event.start_time)
  const camera = findWhere(cameras, { id: event.camera_id })
  if (!camera) return ''

  let duration = `${event.duration}s`
  if (event.duration > 60) {
    duration = timeString(event.end_time)
  }

  return `
    <div class="event col-lg-6 col-lg-offset-3 col-md-8 col-md-offset-2 col-sm-10 col-sm-offset-1" data-id="${event.id}">
      <a href="/event/${event.id}">
        <img class="event-image" src="/video/${event.first_frame}">
        <div class="event-video-progress"></div>
        <div class="event-info trans-opacity">
          <h3>${camera.name}</h3>
          ${time}<br>
          ${duration}<br>
        </div>
      </a>
    </div>
  `
}
