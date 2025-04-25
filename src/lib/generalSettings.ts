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

export function getBackupFreq() {
  return localStorage.getItem("backup-freq") || "every week";
}

export function setBackupFreq(value: string) {
  localStorage.setItem("backup-freq", value);
}
