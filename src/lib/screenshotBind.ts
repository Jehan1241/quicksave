const defaultBind = "-";

export function getScreenshotBind() {
  return localStorage.getItem("screenshot-bind") || defaultBind;
}

export function setScreenshotBind(bind: string) {
  localStorage.setItem("screenshot-bind", bind);
}
