const newCamera = {
  id: 0,
  name: '',
  url: '',
  username: '',
  password: '',
  min_change: '',
  threshold: ''
}

export default function (camera = newCamera) {
  let template = `
    <div class="camera form-horizontal col-sm-8 col-sm-offset-2" data-id="${camera.id}">
      <div class="delete-camera" data-id="${camera.id}">X</div>
      <div class="form-group">
        <label class="control-label col-sm-3">
          Camera Name
        </label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-name form-control"
                 placeholder="Camera Name"
                 value="${camera.name}">
        </div>
      </div>
      <div class="form-group">
        <label class="control-label col-sm-3">URL</label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-url form-control"
                 placeholder="URL"
                 value="${camera.url}">
        </div>
      </div>
      <div class="form-group">
        <label class="control-label col-sm-3">Username</label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-user form-control"
                 placeholder="Username"
                 value="${camera.username}">
        </div>
      </div>
      <div class="form-group">
        <label class="control-label col-sm-3">Password</label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-user form-control"
                 placeholder="Password"
                 value="${camera.password}">
        </div>
      </div>
      <div class="form-group">
        <label class="control-label col-sm-3">Min Change</label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-min-change form-control"
                 placeholder="Min Change"
                 value="${camera.min_change}">
        </div>
      </div>
      <div class="form-group">
        <label class="control-label col-sm-3">Threshold</label>
        <div class="col-sm-8">
          <input type="text"
                 class="camera-threshold form-control"
                 placeholder="Threshold"
                 value="${camera.threshold}">
        </div>
      </div>
  `
  if (camera.id) {
    template += `
    <div class="form-group">
      <label class="control-label col-sm-3">Frame</label>
      <div class="col-sm-8">
        <img class="camera-frame" src="/frame/${camera.id}">
      </div>
    </div>
    `
  }
  template += `
    <div class="camera-buttons" class="btn-group">
      <button type="button" class="btn btn-primary">Save</button>
      <button type="button" class="btn btn-default">Cancel</button>
    </div>
  </div>
  `
  return template
}
