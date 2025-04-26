import { app, BrowserWindow, dialog, ipcMain, shell } from "electron";
import path from "path";
import fs from "fs-extra";
import { spawn, exec } from "child_process";
import extract from "extract-zip";
import { version } from "os";

const REPO = "Jehan1241/quicksave";
const BACKUP_DIR = "backend-backup";

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
    const zipAsset = release.assets.find(
      (a: any) => a.name.startsWith("quicksave-") && a.name.endsWith(".zip")
    );
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
  const appDir = String(process.env.PORTABLE_EXECUTABLE_DIR);
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
    await extract(tempZip, { dir: tempExtract });

    // 3. Verify critical files
    const requiredFiles = [
      "quicksave.exe",
      "backend/thismodule.exe,",
      "backend/updater.exe",
    ];

    for (const file of requiredFiles) {
      if (!fs.existsSync(path.join(tempExtract, file))) {
        throw new Error(`Update package missing required file: ${file}`);
      }
    }

    //4. update the updater
    const extractedUpdaterPath = path.join(
      tempExtract,
      "backend",
      "updater.exe"
    );
    const oldUpdaterPath = path.join(
      String(process.env.PORTABLE_EXECUTABLE_DIR),
      "backend",
      "updater.exe"
    );

    await fs.ensureDir(path.dirname(oldUpdaterPath));
    await fs.copy(extractedUpdaterPath, oldUpdaterPath, { overwrite: true });

    try {
      const response = await fetch("http://localhost:8080/updateApp", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          source: tempExtract,
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
