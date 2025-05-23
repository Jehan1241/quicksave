import { ipcRenderer, contextBridge } from "electron";
import path from "path";

contextBridge.exposeInMainWorld("windowFunctions", {
  nukeCache: () => ipcRenderer.send("nukeCache"),
  closeApp: () => ipcRenderer.send("closeApp"),
  minimize: () => ipcRenderer.send("minimize"),
  maximize: () => ipcRenderer.send("maximize"),
  updatePlayingGame: (uid: string, screenshotBind: string) =>
    ipcRenderer.send("update-playing-game", uid, screenshotBind),
});

contextBridge.exposeInMainWorld("electron", {
  openFolder: (folderPath: string) =>
    ipcRenderer.invoke("open-folder", folderPath),
  browseFileHandler: () => ipcRenderer.invoke("browseFileHandler"),
  validateGamePath: (gamePath: any) =>
    ipcRenderer.invoke("validate-game-path", gamePath),

  onUpdateAvailable: (
    callback: (update: { version: string; zipUrl: string }) => void
  ) => {
    ipcRenderer.on("update-available", (_, update) => callback(update));
  },
  sendUpdateResponse: (choice: boolean) =>
    ipcRenderer.send("update-response", choice),

  onProgress: (callback: (progress: number) => void) => {
    ipcRenderer.on("download-progress", (_, progress) => callback(progress));
  },
  imageSearch: (query: string, page: number) =>
    ipcRenderer.invoke("image-search", query, page),
  fetchImageBuffer: async (url: string): Promise<string | null> => {
    return await ipcRenderer.invoke("fetch-image-buffer", url);
  },
  updateMinimizeSetting: (value: boolean) =>
    ipcRenderer.invoke("update-minimize-setting", value),
  onRequestMinimizeSetting: (callback: (value: boolean) => void) => {
    ipcRenderer.on("request-minimize-setting", () => {
      const value = localStorage.getItem("minimize-on-quit") === "true";
      callback(value);
    });
  },
  sendMinimizeSetting: (value: boolean) => {
    ipcRenderer.send("send-minimize-setting", value);
  },
});
contextBridge.exposeInMainWorld("appPaths", {
  exePath: async () => ipcRenderer.invoke("get-exe-path"),
});

// --------- Expose some API to the Renderer process ---------
contextBridge.exposeInMainWorld("ipcRenderer", {
  on(...args: Parameters<typeof ipcRenderer.on>) {
    const [channel, listener] = args;
    return ipcRenderer.on(channel, (event, ...args) =>
      listener(event, ...args)
    );
  },
  off(...args: Parameters<typeof ipcRenderer.off>) {
    const [channel, ...omit] = args;
    return ipcRenderer.off(channel, ...omit);
  },
  send(...args: Parameters<typeof ipcRenderer.send>) {
    const [channel, ...omit] = args;
    return ipcRenderer.send(channel, ...omit);
  },
  invoke(...args: Parameters<typeof ipcRenderer.invoke>) {
    const [channel, ...omit] = args;
    return ipcRenderer.invoke(channel, ...omit);
  },

  // You can expose other APTs you need here.
  // ...
});

// --------- Preload scripts loading ---------
function domReady(
  condition: DocumentReadyState[] = ["complete", "interactive"]
) {
  return new Promise((resolve) => {
    if (condition.includes(document.readyState)) {
      resolve(true);
    } else {
      document.addEventListener("readystatechange", () => {
        if (condition.includes(document.readyState)) {
          resolve(true);
        }
      });
    }
  });
}

const safeDOM = {
  append(parent: HTMLElement, child: HTMLElement) {
    if (!Array.from(parent.children).find((e) => e === child)) {
      return parent.appendChild(child);
    }
  },
  remove(parent: HTMLElement, child: HTMLElement) {
    if (Array.from(parent.children).find((e) => e === child)) {
      return parent.removeChild(child);
    }
  },
};

/**
 * https://tobiasahlin.com/spinkit
 * https://connoratherton.com/loaders
 * https://projects.lukehaas.me/css-loaders
 * https://matejkustec.github.io/SpinThatShit
 */
function useLoading() {
  const className = `loaders-css__square-spin`;
  const styleContent = `
@keyframes square-spin {
  25% { transform: perspective(100px) rotateX(180deg) rotateY(0); }
  50% { transform: perspective(100px) rotateX(180deg) rotateY(180deg); }
  75% { transform: perspective(100px) rotateX(0) rotateY(180deg); }
  100% { transform: perspective(100px) rotateX(0) rotateY(0); }
}
.${className} > div {
  animation-fill-mode: both;
  width: 50px;
  height: 50px;
  background: #fff;
  animation: square-spin 3s 0s cubic-bezier(0.09, 0.57, 0.49, 0.9) infinite;
}
.app-loading-wrap {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #282c34;
  z-index: 9;
}
    `;
  const oStyle = document.createElement("style");
  const oDiv = document.createElement("div");

  oStyle.id = "app-loading-style";
  oStyle.innerHTML = styleContent;
  oDiv.className = "app-loading-wrap";
  oDiv.innerHTML = `<div class="${className}"><div></div></div>`;

  return {
    appendLoading() {
      safeDOM.append(document.head, oStyle);
      safeDOM.append(document.body, oDiv);
    },
    removeLoading() {
      safeDOM.remove(document.head, oStyle);
      safeDOM.remove(document.body, oDiv);
    },
  };
}

// ----------------------------------------------------------------------

const { appendLoading, removeLoading } = useLoading();
domReady().then(appendLoading);

window.onmessage = (ev) => {
  ev.data.payload === "removeLoading" && removeLoading();
};

setTimeout(removeLoading, 4999);
