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

export function getLastBackupTime(): number {
  return parseInt(localStorage.getItem("last-backup") || "0");
}

export function setLastBackupTime(ts: number) {
  localStorage.setItem("last-backup", ts.toString());
}

export function shouldBackupNow(freq: string, last: number): boolean {
  const now = Date.now();
  const day = 86400000;

  switch (freq) {
    case "on launch":
      return true;
    case "every 2 days":
      return now - last >= 2 * day;
    case "every week":
      return now - last >= 7 * day;
    default:
      return false;
  }
}
