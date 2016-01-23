import cameraConfigTemplate from '../templates/camera-config'
import rowTemplate from '../templates/row'
import {
  contains,
  del,
  filter,
  first,
  group,
  join,
  map,
  partial,
  pipe
} from '../utils/utils'
import {
  element,
  findClass,
  on,
  prepend,
  remove
} from '../utils/dom'

const camerasDiv = document.getElementById('cameras')
const cameras = window.OPENCAM_CAMERAS
const addCameraBtn = document.getElementById('add-camera-btn')

camerasDiv.innerHTML = rowTemplate(
  join(map(cameras, cameraConfigTemplate))
)

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
  cameraConfigTemplate,
  element,
  first,
  partial(prepend, camerasRow),
  bindCameraEvents
))

function bindCameraEvents (cameraEl) {
  let cameraId = cameraEl.dataset.id
  let deleteEl = findClass(cameraEl, 'delete-camera')
  on(deleteEl, 'click', () => {
    if (cameraId > 0) {
      cameraElements = del(cameraElements, cameraId)
      // ajax delete camera
    }
    remove(cameraEl)
  })
}
