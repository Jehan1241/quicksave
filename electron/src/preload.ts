const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("windowFunctions", {
    closeApp: () => ipcRenderer.send("closeApp"),
    minimize: () => ipcRenderer.send("minimize"),
    maximize: () => ipcRenderer.send("maximize"),
});
