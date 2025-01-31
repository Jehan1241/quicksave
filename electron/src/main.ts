import { app, BrowserWindow, ipcMain } from "electron";
import path from "path";
const ipc = ipcMain;

const inDevelopment = process.env.NODE_ENV === "development";

function createWindow() {
    const preload = path.join(__dirname, "preload.js");
    const mainWindow = new BrowserWindow({
        width: 1080,
        height: 550,
        minWidth: 1080,
        minHeight: 550,
        webPreferences: {
            devTools: inDevelopment,
            contextIsolation: true,
            nodeIntegration: true,
            nodeIntegrationInSubFrames: false,
            preload: preload,
        },
        titleBarStyle: "hidden",
        autoHideMenuBar: true,
    });

    ipc.on("closeApp", () => {
        mainWindow.close();
    });
    ipc.on("minimize", () => {
        mainWindow.minimize();
    });
    ipc.on("maximize", () => {
        if (mainWindow.isMaximized()) {
            mainWindow.restore();
        } else {
            mainWindow.maximize();
        }
    });

    if (MAIN_WINDOW_VITE_DEV_SERVER_URL) {
        mainWindow.loadURL(MAIN_WINDOW_VITE_DEV_SERVER_URL);
    } else {
        mainWindow.loadFile(
            path.join(__dirname, `../renderer/${MAIN_WINDOW_VITE_NAME}/index.html`)
        );
    }
}

app.disableHardwareAcceleration();
app.whenReady().then(createWindow);

//osX only
app.on("window-all-closed", () => {
    if (process.platform !== "darwin") {
        app.quit();
    }
});

app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
        createWindow();
    }
});
//osX only ends
