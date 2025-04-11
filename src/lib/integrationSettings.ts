export function getIntegrateOnLaunchEnabled() {
  const state = localStorage.getItem("integrate-on-launch") || "true";
  if (state === "true") return true;
  else return false;
}

export function setIntegrateOnLaunchEnabled(state: boolean) {
  if (state) localStorage.setItem("integrate-on-launch", "true");
  else localStorage.setItem("integrate-on-launch", "false");
}

export function getIntegrateOnExitEnabled() {
  const state = localStorage.getItem("integrate-on-exit") || "true";
  if (state === "true") return true;
  else return false;
}

export function setIntegrateOnExitEnabled(state: boolean) {
  if (state) localStorage.setItem("integrate-on-exit", "true");
  else localStorage.setItem("integrate-on-exit", "false");
}
