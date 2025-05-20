import { app, BrowserWindow, ipcMain } from "electron";
import path from "path";
import fs from "fs-extra";
import { exec } from "child_process";

import { promisify } from "util";

// Properly declare execAsync
const execAsync = promisify(exec);

const REPO = "Jehan1241/quicksave";

const appPrefix =
  process.platform === "win32" ? "quicksave-windows-" : "quicksave-linux-";

const requiredFiles =
  process.platform === "win32"
    ? [
        "quicksave/quicksave.exe",
        "quicksave/backend/quicksaveService.exe",
        "quicksave/backend/updater.exe",
      ]
    : ["quicksave/quicksave", "quicksaveService", "updater"];

export async function checkForUpdates(currentVersion: any) {
  try {
    // 1. Fetch latest release from GitHub
    const response = await fetch(
      `https://api.github.com/repos/${REPO}/releases/latest`
    );
    const release = await response.json();
    console.log("RELEASE", release.tag_name);

    if (release.tag_name === `v${currentVersion}`) return null;

    // 2. Find the portable ZIP asset
    let zipAsset;
    zipAsset = release.assets.find(
      (a: any) => a.name.startsWith(appPrefix) && a.name.endsWith(".zip")
    );
    if (!zipAsset) {
      zipAsset = release.assets.find(
        //for legacy support
        (a: any) => a.name.startsWith("quicksave-") && a.name.endsWith(".zip")
      );
    }
    if (!zipAsset) throw new Error("No portable ZIP found in release");

    console.log("RELEASE", zipAsset.browser_download_url);
    return {
      version: release.tag_name,
      zipUrl: zipAsset.browser_download_url,
    };
  } catch (error) {
    console.error("Update check failed:", error);
    return null;
  }
}

async function downloadUpdate(zipUrl: string, win: BrowserWindow) {
  const appDir = path.dirname(process.execPath);
  const tempZip = path.join(app.getPath("temp"), "quicksave-update.zip");
  console.log(`Temp ZIP path: ${tempZip}`);
  const tempExtract = path.join(app.getPath("temp"), "quicksave-extracted");

  try {
    console.log(`Downloading update...`);
    const response = await fetch(zipUrl);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    if (!response.body) throw new Error(`Response body is null`);

    const contentLength = Number(response.headers.get("content-length")) || 0;
    console.log(`Expected size: ${contentLength} bytes`);

    const reader = response.body.getReader();
    let receivedLength = 0;
    const chunks = [];
    let lastSentPercent = -1;

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      chunks.push(value);
      receivedLength += value.length;

      if (contentLength > 0) {
        const currentPercent = Math.round(
          (receivedLength / contentLength) * 100
        );
        if (currentPercent !== lastSentPercent) {
          win.webContents.send("download-progress", currentPercent);
          lastSentPercent = currentPercent;
        }
      }
    }

    win.webContents.send("download-progress", 100);

    // 2. Combine chunks using modern Blob API
    const blob = new Blob(chunks);
    const arrayBuffer = await blob.arrayBuffer();
    await fs.writeFile(tempZip, Buffer.from(arrayBuffer));

    // 2. Extract to temp directory first
    await fs.ensureDir(tempExtract);
    try {
      if (process.platform === "win32") {
        await execAsync(
          `powershell Expand-Archive -Path "${tempZip}" -DestinationPath "${tempExtract}" -Force`
        );
      } else {
        await execAsync(`bsdtar -xf "${tempZip}" -C "${tempExtract}"`);
      }
    } catch (error) {
      console.log("Extraction Error: ", error);
    }

    // 3. Verify critical files
    for (const file of requiredFiles) {
      if (!fs.existsSync(path.join(tempExtract, file))) {
        throw new Error(`Update package missing required file: ${file}`);
      }
    }

    //4. update the updater
    const extractedUpdaterPath = path.join(
      tempExtract,
      "quicksave",
      "backend",
      process.platform === "win32" ? "updater.exe" : "updater"
    );

    const oldUpdaterPath = path.join(
      path.dirname(process.execPath),
      "backend",
      process.platform === "win32" ? "updater.exe" : "updater"
    );
    await fs.ensureDir(path.dirname(oldUpdaterPath));
    await fs.copy(extractedUpdaterPath, oldUpdaterPath, { overwrite: true });

    const source = path.join(tempExtract, "quicksave");

    try {
      const response = await fetch("http://localhost:50001/updateApp", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          source: source,
          target: appDir,
        }),
      });
      if (!response.ok) throw new Error("Update initiation failed");
    } catch (error) {
      console.error("Update failed:", error);
    }
  } catch (error) {
    console.error("Update failed:", error);
    throw error;
  } finally {
    await fs.remove(tempZip).catch(() => {});
  }
}

export async function promptUpdate(win: BrowserWindow, currentVersion: any) {
  const update = await checkForUpdates(currentVersion);
  if (!update) return false;

  win.webContents.once("did-finish-load", () => {
    setTimeout(() => {
      win.webContents.send("update-available", {
        version: update.version,
        zipUrl: update.zipUrl,
      });
    }, 2000); //Avoid edge case no dialog showing
  });

  const response = await new Promise<boolean>((resolve) => {
    ipcMain.once("update-response", (_, userChoice) => resolve(userChoice));
  });

  if (response) {
    await downloadUpdate(update.zipUrl, win);
  }
}
