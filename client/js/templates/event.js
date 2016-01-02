import moment from 'moment'

export default function (event) {
  const time = moment(event.start_time).format('hh:mm:ssa MM/DD/YYYY')
  return `
    <div class="event col-lg-3 col-md-4 col-sm-6" data-id="${event.id}">
      <a href="/event/${event.id}">
        <div class="event-image-container">
          <img class="event-image" src="/video/${event.first_frame}">
        </div>
        <div class="event-info">
          ${time}<br>
          ${event.duration}s<br>
        </div>
      </a>
    </div>
  `
}
