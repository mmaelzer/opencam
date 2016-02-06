import moment from 'moment'
import { findWhere } from '../utils/utils'

export default function (event, cameras) {
  const time = moment(event.start_time).format('hh:mm:ssa MM/DD/YYYY')
  const camera = findWhere(cameras, { id: event.camera_id })
  if (!camera) return ''
  return `
    <div class="event fadein col-lg-3 col-md-4 col-sm-6" data-id="${event.id}">
      <a href="/event/${event.id}">
        <img class="event-image" src="/video/${event.first_frame}">
        <div class="event-video-progress"></div>
        <div class="event-info trans-opacity">
          <h3>${camera.name}</h3>
          ${time}<br>
          ${event.duration}s<br>
        </div>
      </a>
    </div>
  `
}
