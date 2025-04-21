import { app, BrowserWindow, shell, ipcMain, dialog } from "electron";
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import path from "node:path";
import { promptUpdate } from "./github-updater";
const require = createRequire(import.meta.url);
import { Tray, nativeImage, Menu } from "electron";
const { globalShortcut } = require("electron");

process.env.ELECTRON_ENABLE_LOGGING = "1";
console.log("=== MAIN PROCESS STARTED ===");

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
const windowStateKeeper = require("electron-window-state");
let goServer: any;
const { spawn } = require("child_process");
const isDev = !app.isPackaged;
const fs = require("fs");

let tray: Tray | null = null;

async function createWindow() {
  let mainWindowState = windowStateKeeper({
    defaultWidth: 1080,
    defaultHeight: 550,
  });

  const win = new BrowserWindow({
    x: mainWindowState.x,
    y: mainWindowState.y,
    width: mainWindowState.width,
    height: mainWindowState.height,
    minWidth: 1080,
    minHeight: 550,
    title: "quicksave",
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

  mainWindowState.manage(win);

  const { ipcMain } = require("electron");
  const fs = require("fs");
  const path = require("path");
  const ws = require("windows-shortcuts");

  ipcMain.handle(
    "browseFileHandler",
    async (_: any, options: Electron.OpenDialogOptions) => {
      return await dialog.showOpenDialog(options);
    }
  );

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

  let minimizeToTray = true;

  // Add this to initialize when the window is created
  win.webContents.on("did-finish-load", () => {
    win.webContents.send("request-minimize-setting");
  });

  ipcMain.on("send-minimize-setting", (event: any, value: boolean) => {
    minimizeToTray = value;
    console.log("Minimize to Tray initialized:", value);
  });

  ipcMain.handle("update-minimize-setting", (event: any, value: any) => {
    minimizeToTray = value;
  });

  ipc.on("closeApp", () => {
    if (minimizeToTray) win.hide();
    else win.close();
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

  ipcMain.on(
    "update-playing-game",
    (_: any, uid: string, screenshotBind: string) => {
      unregisterScreenshotShortcut(screenshotBind);
      if (uid === "") {
        unregisterScreenshotShortcut(screenshotBind);
      } else if (uid) {
        registerScreenshotShortcut(uid, screenshotBind);
      }
    }
  );

  ipcMain.handle("open-folder", (event: any, folderPath: string) => {
    const exePath = isDev
      ? path.resolve(__dirname, "../..")
      : String(process.env.PORTABLE_EXECUTABLE_DIR);
    shell
      .openPath(path.join(exePath, folderPath))
      .catch((err) => console.error("Error opening folder:", err));
  });

  if (VITE_DEV_SERVER_URL) {
    win.loadURL(VITE_DEV_SERVER_URL);
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
  if (app.isPackaged) promptUpdate(win, app.getVersion());

  const iconPath = path.join(__dirname, "/../../assets/image.png");
  const trayIcon = nativeImage.createFromPath(iconPath);

  tray = new Tray(trayIcon);
  tray.setToolTip("quicksave");

  const contextMenu = Menu.buildFromTemplate([
    {
      label: "Show App",
      click: () => win.show(),
    },
    {
      label: "Quit",
      click: () => {
        tray?.destroy();
        app.quit();
      },
    },
  ]);

  // Assign context menu to tray
  tray.setContextMenu(contextMenu);

  // Left click to show (existing code)
  tray.on("click", () => {
    win.show();
  });
}

app.whenReady().then(() => {
  createWindow();

  app.on("will-quit", () => {
    globalShortcut.unregisterAll(); // Clean up when app exits
  });
});

function ensureBackend() {
  const exeDest = path.join(
    String(process.env.PORTABLE_EXECUTABLE_DIR),
    "backend",
    "thismodule.exe"
  );
  if (!fs.existsSync(exeDest)) {
    console.error("Go server (thismodule.exe) not found");
  }
  return exeDest;
}

app.on("ready", () => {
  let serverPath;

  if (isDev) {
    serverPath = path.join(__dirname, "../../backend", "thismodule.exe");
  } else {
    serverPath = ensureBackend();
  }
  console.log("Launching Go server from:", serverPath);

  goServer = spawn(serverPath, [], {
    cwd: path.dirname(serverPath),
    shell: false,
    stdio: ["ignore", "pipe", "pipe"], // Capture stdout and stderr
  });

  goServer.stdout.on("data", (data: any) => {
    const msg = `[STDOUT] ${data.toString()}`;
    console.log(msg);
  });

  goServer.stderr.on("data", (data: any) => {
    const msg = `[STDERR] ${data.toString()}`;
    console.error(msg);
  });

  goServer.on("error", (err: any) => {
    const msg = `[ERROR] Failed to start Go server: ${err.message}`;
    console.error(msg);
  });

  goServer.on("exit", (code: any) => {
    const msg = `[EXIT] Go server exited with code ${code}`;
    console.log(msg);
  });
});

const killGoServer = () => {
  if (goServer && !goServer.killed) {
    console.log("Attempting to kill Go server...");
    goServer.kill("SIGTERM");

    setTimeout(() => {
      if (!goServer.killed) {
        console.log("Force killing Go server...");

        if (process.platform === "win32") {
          require("child_process").exec(`taskkill /PID ${goServer.pid} /F`);
        } else {
          process.kill(-goServer.pid, "SIGKILL"); // Kill entire process group
        }
      }
    }, 2000);
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
  console.log("Here");
  goServer.kill();
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

function registerScreenshotShortcut(uid: string, screenshotBind: string) {
  globalShortcut.register(screenshotBind, async () => {
    if (!uid) {
      console.log("No game UID found, skipping screenshot request.");
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/takeScreenshot?uid=${uid}`
      );
      const data = await response.text();
      console.log("Screenshot request sent, response:", data);
    } catch (error) {
      console.error("Error sending screenshot request:", error);
    }
  });
}

function unregisterScreenshotShortcut(screenshotBind: string) {
  if (globalShortcut.isRegistered(screenshotBind)) {
    globalShortcut.unregister(screenshotBind);
  }
}

ipcMain.handle("image-search", async (_event, query, page) => {
  const userAgents = [
    // Chrome (Windows/Mac/Linux)
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",

    // Firefox (Windows/Mac/Linux)
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 13.5; rv:121.0) Gecko/20100101 Firefox/121.0",
    "Mozilla/5.0 (X11; Linux i686; rv:119.0) Gecko/20100101 Firefox/119.0",

    // Safari (Mac/iOS)
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",

    // Edge (Windows/Mac)
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.2210.91",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.2210.91",

    // Mobile Devices
    "Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.210 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; SM-A736B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.210 Mobile Safari/537.36",

    // Additional Variants
    "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko", // IE11
    "Mozilla/5.0 (Windows NT 6.3; Win64; x64; Trident/7.0; Touch; rv:11.0) like Gecko", // IE11 Touch
    "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", // ChromeOS
  ];

  // Improved random user agent selection
  const randomUserAgent =
    userAgents[Math.floor(Math.random() * userAgents.length)];

  // Additional headers for better emulation
  const headers = {
    Accept:
      "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    Connection: "keep-alive",
    "Upgrade-Insecure-Requests": "1",
    "Sec-Fetch-Dest": "document",
    "Sec-Fetch-Mode": "navigate",
    "Sec-Fetch-Site": "none",
    "Sec-Fetch-User": "?1",
    "Cache-Control": "max-age=0",
  };

  const win = new BrowserWindow({
    show: false,
    webPreferences: {
      offscreen: true,
      contextIsolation: true,
      nodeIntegration: false,
      webSecurity: true,
    },
  });

  // Set headers and user agent
  await win.webContents.setUserAgent(randomUserAgent);
  await win.webContents.session.webRequest.onBeforeSendHeaders(
    (details, callback) => {
      details.requestHeaders = { ...details.requestHeaders, ...headers };
      callback({ requestHeaders: details.requestHeaders });
    }
  );

  const url = `https://www.google.com/search?hl=en&tbm=isch&q=${encodeURIComponent(query)}&start=${page}`;
  console.log("Loading URL:", url);

  try {
    // Load page with timeout
    await Promise.race([
      win.loadURL(url),
      new Promise((_, reject) =>
        setTimeout(() => reject(new Error("Page load timeout")), 15000)
      ),
    ]);

    await win.webContents.executeJavaScript(`
      window.scrollBy(0, 100); // Scroll down by 100px
      setTimeout(() => {}, 500); // Wait a bit before loading more content
    `);

    // Handle consent form if it appears (like Playnite does)
    const hasConsent = await win.webContents.executeJavaScript(`
          document.querySelector('form[action*="consent.google.com"]') !== null
        `);

    if (hasConsent) {
      await win.webContents.executeJavaScript(`
            document.querySelector('form').submit();
          `);
      await new Promise((resolve) => setTimeout(resolve, 2000));
    }

    // Wait for either metadata or images to appear (with timeout)
    await win.webContents.executeJavaScript(`
      new Promise((resolve, reject) => {
        const timeout = setTimeout(() => reject(new Error('Content load timeout')), 10000);

        const checkReady = () => {
          if (document.querySelector('.rg_meta') || document.querySelector('img[src^="http"]')) {
            clearTimeout(timeout);
            resolve();
          } else {
            setTimeout(checkReady, 300);
          }
        };

        checkReady();
      });
    `);

    const images = await win.webContents.executeJavaScript(`
      (() => {
        try {
          const results = [];
          const decodeUrl = (url) => {
            try {
              return decodeURIComponent(url);
            } catch (e) {
              return url; // Return original if decoding fails
            }
          };

          // First: Try Google's classic .rg_meta style
          const metaElements = document.querySelectorAll('.rg_meta');
          if (metaElements.length > 0) {
            metaElements.forEach(meta => {
              try {
                const data = JSON.parse(meta.textContent);
                const imageUrl = data.ou || data.ru;
                const thumbUrl = data.tu;

                results.push({
                  ImageUrl: decodeUrl(imageUrl),
                  ThumbUrl: decodeUrl(thumbUrl),
                  Width: data.ow,
                  Height: data.oh,
                  position: parseInt(meta.closest('[data-ri]')?.getAttribute('data-ri') || 0)
                });
              } catch (e) {
                console.warn('Metadata parse error:', e);
              }
            });
            return results;
          }

          // Fallback: Regex scrape like Playnite
          const html = document.documentElement.outerHTML.replace(/\\n/g, "");
          const regex = /\\["(https:\\/\\/encrypted-[^,]+?)",\\d+,\\d+\\],\\["(http.+?)",(\\d+),(\\d+)\\]/g;
          let match;
          while ((match = regex.exec(html)) !== null) {
            results.push({
              ThumbUrl: decodeUrl(match[1]),
              ImageUrl: decodeUrl(match[2]),
              Width: parseInt(match[4]),
              Height: parseInt(match[3]),
              position: results.length
            });
          }

          return results;
        } catch (e) {
          console.error('Image scraping failed:', e);
          return [];
        }
      })();
    `);

    // Log all image URLs found
    console.log("Raw image URLs scraped:");
    images.forEach((img: any) => {
      console.log(`Raw Image URL: ${img.ImageUrl}`);
    });

    // Sort by original position (like Google's ranking)
    images.sort((a: any, b: any) => a.position - b.position);

    // Filter to only valid, high-res images (minimum 300px on either dimension)
    const validImages = images.filter(
      (img: any) =>
        img.ImageUrl?.startsWith("http") &&
        img.ThumbUrl?.startsWith("http") &&
        Math.max(img.Width, img.Height) >= 300
    );

    // Log filtered (valid) image URLs
    console.log("Filtered valid image URLs:");
    validImages.forEach((img: any) => {
      console.log(`Valid Image URL: ${img.ImageUrl}`);
    });

    // Return top 20 results (like Google's pagination)
    const topImages = validImages;
    console.log(`Returning ${topImages.length} valid image results`);
    win.destroy();
    return topImages;
  } catch (err) {
    console.error("Image search failed:", err);
    if (!win.isDestroyed()) win.destroy();
    return []; // Return empty array instead of throwing error
  }
});
