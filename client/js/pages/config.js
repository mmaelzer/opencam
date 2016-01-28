import cameraConfigTemplate from '../templates/camera-config'
import rowTemplate from '../templates/row'
import {
  breaker,
  captureErr,
  contains,
  del,
  filter,
  group,
  join,
  map,
  partial,
  pipe,
  reduce
} from '../utils/utils'
import {
  element,
  findClass,
  findTags,
  on,
  prepend,
  remove,
  replace
} from '../utils/dom'
import superagent from 'superagent'
import superagentJson from 'superagent-jsonapify'
superagentJson(superagent)

const request = superagent
const camerasDiv = document.getElementById('cameras')
const cameras = window.OPENCAM_CAMERAS
const camerasById = group(window.OPENCAM_CAMERAS, 'id')
const addCameraBtn = document.getElementById('add-camera-btn')

camerasDiv.innerHTML = rowTemplate(
  join(map(cameras, cameraConfigTemplate))
)

const reload = () => window.location.href = window.location.href

const camerasRow = camerasDiv.firstChild
map(
  filter(
    camerasRow.childNodes, (node) => contains(node.className, 'camera')
  ),
  bindCameraEvents
)

let cameraElements = group(
  document.getElementsByClassName('camera'), (el) => el.dataset.id
)

on(addCameraBtn, 'click', pipe(
  partial(breaker, () => 0 in camerasById),
  cameraConfigTemplate,
  element,
  partial(prepend, camerasRow),
  bindCameraEvents,
  () => camerasById[0] = {}
))

function gatherCameraInfo (el) {
  let inputs = findTags(el, 'input')
  return reduce(inputs, (cam, input) => {
    cam[input.name] = isNaN(+input.value) ? input.value : +input.value
    return cam
  }, {})
}

function bindCameraEvents (cameraEl) {
  let cameraId = cameraEl.dataset.id
  let deleteBtn = findClass(cameraEl, 'delete-camera')
  let saveBtn = findClass(cameraEl, 'save-btn')
  let cancelBtn = findClass(cameraEl, 'cancel-btn')
  on(saveBtn, 'click', () => {
    let info = gatherCameraInfo(cameraEl)
    request.post('/api/camera')
           .send(info)
           .end(captureErr(reload))
  })

  on(cancelBtn, 'click', () => {
    bindCameraEvents(
      replace(cameraEl,
        element(
          cameraConfigTemplate(camerasById[cameraId])
        )
      )
    )
  })

  on(deleteBtn, 'click', () => {
    if (cameraId > 0) {
      // ajax delete camera
    }
    cameraElements = del(cameraElements, cameraId)
    remove(cameraEl)
  })

  return cameraEl
}
