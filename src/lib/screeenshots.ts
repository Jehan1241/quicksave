import { useSortContext } from "@/hooks/useSortContex";

const defaultBind = "-";

export function getScreenshotBind() {
  return localStorage.getItem("screenshot-bind") || defaultBind;
}

export function setScreenshotBind(bind: string) {
  localStorage.setItem("screenshot-bind", bind);
}

export function getScreenshotEnabled() {
  const state = localStorage.getItem("screenshot-enabled") || "enabled";
  if (state === "enabled") return true;
  else return false;
}

export function setScreenshotEnabled(state: boolean) {
  if (state) localStorage.setItem("screenshot-enabled", "enabled");
  else localStorage.setItem("screenshot-enabled", "disabled");
}
