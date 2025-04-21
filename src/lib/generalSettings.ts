const minimizeOnQuitDefault = "true";

export function getMinimizedToTray() {
  const value =
    localStorage.getItem("minimize-on-quit") || minimizeOnQuitDefault;
  if (value === "true") return true;
  else return false;
}

export async function setMinimizeToTray(value: boolean) {
  localStorage.setItem("minimize-on-quit", value ? "true" : "false");
  window.electron.updateMinimizeSetting(value);
}
