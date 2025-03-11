import { app, BrowserWindow, shell, ipcMain, session } from "electron";
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import path, { dirname } from "node:path";
import os from "node:os";
import { update } from "./update";

const require = createRequire(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ipc = ipcMain;

// The built directory structure
//
// ├─┬ dist-electron
// │ ├─┬ main
// │ │ └── index.js    > Electron-Main
// │ └─┬ preload
// │   └── index.mjs   > Preload-Scripts
// ├─┬ dist
// │ └── index.html    > Electron-Renderer
//
process.env.APP_ROOT = path.join(__dirname, "../..");

export const MAIN_DIST = path.join(process.env.APP_ROOT, "dist-electron");
export const RENDERER_DIST = path.join(process.env.APP_ROOT, "dist");
export const VITE_DEV_SERVER_URL = process.env.VITE_DEV_SERVER_URL;

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL
  ? path.join(process.env.APP_ROOT, "public")
  : RENDERER_DIST;

// // Disable GPU Acceleration for Windows 7
// if (os.release().startsWith("6.1")) app.disableHardwareAcceleration();

// // Set application name for Windows 10+ notifications
// if (process.platform === "win32") app.setAppUserModelId(app.getName());

// if (!app.requestSingleInstanceLock()) {
//   app.quit();
//   process.exit(0);
// }

let win: BrowserWindow | null = null;
const preload = path.join(__dirname, "../preload/index.mjs");
const indexHtml = path.join(RENDERER_DIST, "index.html");

async function createWindow() {
  const win = new BrowserWindow({
    width: 1080,
    height: 550,
    minWidth: 1080,
    minHeight: 550,
    webPreferences: {
      devTools: true,
      contextIsolation: true,
      nodeIntegration: false,
      nodeIntegrationInSubFrames: false,
      preload,
      webSecurity: false,
    },
    titleBarStyle: "hidden",
    autoHideMenuBar: true,
  });

  const { ipcMain, app } = require("electron");
  const fs = require("fs");
  const path = require("path");
  const ws = require("windows-shortcuts"); // Import windows-shortcuts

  ipcMain.handle("validate-game-path", async (event: any, gamePath: any) => {
    // Step 1: Trim any leading/trailing whitespace from the path
    gamePath = gamePath.trim();

    // Step 2: Remove surrounding quotes (if any) from the path
    if (gamePath.startsWith('"') && gamePath.endsWith('"')) {
      gamePath = gamePath.slice(1, -1); // Remove the leading and trailing quotes
    }

    // Step 3: Check if path includes %USERPROFILE% or Public Desktop
    const userProfile = process.env.USERPROFILE; // Gets the value of %USERPROFILE%
    const publicDesktop = path.join("C:\\Users", "Public", "Desktop");

    if (gamePath.includes("%USERPROFILE%\\Desktop")) {
      gamePath = gamePath.replace(
        "%USERPROFILE%\\Desktop",
        path.join(userProfile, "Desktop")
      );
    } else if (gamePath.includes("%PUBLIC%\\Desktop")) {
      gamePath = gamePath.replace("%PUBLIC%\\Desktop", publicDesktop);
    }

    console.log("Resolved Path:", gamePath); // For debugging purposes

    // Step 4: Check if the path is a .lnk file (Windows shortcut)
    if (path.extname(gamePath).toLowerCase() === ".lnk") {
      try {
        // Return a Promise to handle async behavior properly
        return new Promise((resolve, reject) => {
          ws.query(gamePath, (error: any, options: any) => {
            if (error) {
              console.error("Error reading shortcut:", error);
              reject({
                isValid: false,
                message: "Error reading shortcut target path.",
              });
            }

            // Log the shortcut information
            console.log("Shortcut Information:", options);

            // Extract the expanded target path
            const targetPath = options.expanded.target;

            // Log the resolved executable path for debugging
            console.log("Resolved Executable Path:", targetPath);

            // Check if the executable exists
            if (fs.existsSync(targetPath)) {
              const fileExtension = path.extname(targetPath).toLowerCase();
              const validExtensions = [".exe", ".bin", ".app", ".sh", ".jar"];

              if (validExtensions.includes(fileExtension)) {
                console.log("In here");
                resolve({ isValid: true, message: targetPath }); // Valid game path with exe path from shortcut
              } else {
                reject({
                  isValid: false,
                  message: "Invalid file extension. (Do not link Shortcuts)",
                });
              }
            } else {
              reject({
                isValid: false,
                message: "Target path from shortcut does not exist.",
              });
            }
          });
        });
      } catch (error) {
        // Catch any errors not handled by the callback
        console.error("Error querying shortcut:", error);
        return {
          isValid: false,
          message: "Error reading shortcut target path.",
        };
      }
    } else {
      // Step 5: If it's not a shortcut, proceed with the regular file validation
      if (fs.existsSync(gamePath)) {
        const fileExtension = path.extname(gamePath).toLowerCase();
        const validExtensions = [".exe", ".bin", ".app", ".sh", ".jar"];

        if (validExtensions.includes(fileExtension)) {
          return { isValid: true, message: gamePath }; // Valid game path
        } else {
          return { isValid: false, message: "Invalid file extension." };
        }
      } else {
        return { isValid: false, message: "Path does not exist." };
      }
    }
  });

  // ipc.on("nukeCache", () => {
  //   session.defaultSession.clearCache();
  //   win.webContents.reloadIgnoringCache();
  //   win.webContents.session.clearCache().then(() => {
  //     win.webContents.executeJavaScript("location.reload(true)");
  //   });
  //   win.webContents.invalidate();
  //   win.webContents.reloadIgnoringCache();
  //   win.webContents.session.clearCache();
  //   win.webContents.session.flushStorageData();
  // });

  ipc.on("closeApp", () => {
    win.close();
  });
  ipc.on("minimize", () => {
    win.minimize();
  });
  ipc.on("maximize", () => {
    if (win.isMaximized()) {
      win.restore();
    } else {
      win.maximize();
    }
  });

  if (VITE_DEV_SERVER_URL) {
    // #298
    win.loadURL(VITE_DEV_SERVER_URL);
    // Open devTool if the app is not packaged
    win.webContents.openDevTools();
  } else {
    win.loadFile(indexHtml);
  }

  // Test actively push message to the Electron-Renderer
  win.webContents.on("did-finish-load", () => {
    win?.webContents.send("main-process-message", new Date().toLocaleString());
  });

  // Make all links open with the browser, not with the application
  win.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url);
    return { action: "deny" };
  });
  // This will intercept anchor `<a>` tag clicks
  win.webContents.on("will-navigate", (event, url) => {
    if (!url.startsWith("file://")) {
      event.preventDefault();
      shell.openExternal(url);
    }
  });

  // Auto update
  update(win);
}

app.whenReady().then(createWindow);

let goServer: any;
const { spawn } = require("child_process");

const isDev = !app.isPackaged;
const fs = require("fs");
const logPath = path.join(process.cwd(), "server.log");
const logStream = fs.createWriteStream(logPath, { flags: "a" });

app.on("ready", () => {
  let serverPath;

  if (isDev) {
    // Development: Look for backend in the project directory
    serverPath = path.join(__dirname, "../../backend", "thismodule.exe");
  } else {
    // Production: The backend folder is next to the packaged Electron executable
    const execDir = path.dirname(process.execPath); // Use execPath instead of __dirname
    serverPath = path.join(execDir, "backend", "thismodule.exe");
  }
  console.log("Launching Go server from:", serverPath);

  const goServer = spawn(serverPath, [], {
    cwd: path.dirname(serverPath),
    shell: false,
    stdio: ["ignore", "pipe", "pipe"], // Capture stdout and stderr
  });

  goServer.stdout.on("data", (data: any) => {
    const msg = `[STDOUT] ${data.toString()}`;
    console.log(msg);
    logStream.write(msg);
  });

  goServer.stderr.on("data", (data: any) => {
    const msg = `[STDERR] ${data.toString()}`;
    console.error(msg);
    logStream.write(msg);
  });

  goServer.on("error", (err: any) => {
    const msg = `[ERROR] Failed to start Go server: ${err.message}`;
    console.error(msg);
    logStream.write(msg);
  });

  goServer.on("exit", (code: any) => {
    const msg = `[EXIT] Go server exited with code ${code}`;
    console.log(msg);
    logStream.write(msg);
  });
});

const killGoServer = () => {
  if (goServer && !goServer.killed) {
    console.log("Attempting to kill Go server...");
    goServer.kill("SIGTERM"); // Graceful shutdown

    setTimeout(() => {
      if (!goServer.killed) {
        console.log("Force killing Go server...");
        goServer.kill("SIGKILL"); // Hard kill if needed
      }
    }, 2000); // Wait 2s before force kill
  }
};

app.on("before-quit", killGoServer);
app.on("window-all-closed", killGoServer);
app.on("quit", killGoServer);
process.on("exit", killGoServer);
process.on("SIGINT", killGoServer);
process.on("SIGTERM", killGoServer);

app.on("window-all-closed", () => {
  console.log("And");
  if (goServer) {
    console.log("Here");
    goServer.kill();
  }
  win = null;
  if (process.platform !== "darwin") app.quit();
});

app.on("second-instance", () => {
  if (win) {
    // Focus on the main window if the user tried to open another
    if (win.isMinimized()) win.restore();
    win.focus();
  }
});

app.on("activate", () => {
  const allWindows = BrowserWindow.getAllWindows();
  if (allWindows.length) {
    allWindows[0].focus();
  } else {
    createWindow();
  }
});

// New window example arg: new windows url
ipcMain.handle("open-win", (_, arg) => {
  const childWindow = new BrowserWindow({
    webPreferences: {
      preload,
      nodeIntegration: true,
      contextIsolation: false,
    },
  });

  if (VITE_DEV_SERVER_URL) {
    childWindow.loadURL(`${VITE_DEV_SERVER_URL}#${arg}`);
  } else {
    childWindow.loadFile(indexHtml, { hash: arg });
  }
});
