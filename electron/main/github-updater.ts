import { app, dialog, shell } from "electron";
import path from "path";
import fs from "fs-extra";
import { spawn, exec } from "child_process";
import extract from "extract-zip";

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

async function downloadUpdate(zipUrl: string) {
  const appDir = String(process.env.PORTABLE_EXECUTABLE_DIR);
  const tempZip = path.join(app.getPath("temp"), "quicksave-update.zip");
  const tempExtract = path.join(app.getPath("temp"), "quicksave-extracted");

  try {
    // 1. Download the update
    console.log(`Downloading update...`);
    const response = await fetch(zipUrl);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    await fs.writeFile(tempZip, Buffer.from(await response.arrayBuffer()));

    // 2. Extract to temp directory first
    await fs.ensureDir(tempExtract);
    await extract(tempZip, { dir: tempExtract });

    // 3. Verify critical files
    const requiredFiles = ["quicksave.exe", "backend/thismodule.exe"];

    for (const file of requiredFiles) {
      if (!fs.existsSync(path.join(tempExtract, file))) {
        throw new Error(`Update package missing required file: ${file}`);
      }
    }

    // 4. Prepare update instructions for next launch
    const pendingUpdate = {
      source: tempExtract,
      target: appDir,
      files: requiredFiles,
      timestamp: Date.now(),
    };

    console.log("Pending Update", pendingUpdate);

    try {
      const response = await fetch("http://localhost:8080/updateApp", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          source: pendingUpdate.source,
          target: appDir,
        }),
      });

      if (!response.ok) throw new Error("Update initiation failed");
    } catch (error) {
      console.error("Update failed:", error);
    }
    console.log("After Pending Update", pendingUpdate);

    return true;
  } catch (error) {
    console.error("Update failed:", error);
    throw error;
  } finally {
    // Cleanup temp files (keep extracted files if pending update exists)
    await fs.remove(tempZip).catch(() => {});
  }
}

export async function promptUpdate(currentVersion: any): Promise<boolean> {
  const update = await checkForUpdates(currentVersion);
  if (!update) return false;

  const dialogResult = await dialog.showMessageBox({
    type: "info",
    buttons: ["Update", "Release Notes", "Later"],
    message: `Update to ${update.version} available`,
    detail: "Would you like to download and install it now?",
    cancelId: 2, // Index of "Later" button
    noLink: true, // Prevent default Escape behavior
  });

  // response will be:
  // 0 - Update
  // 1 - Release Notes
  // 2 - Later (or close/X button)
  const { response } = dialogResult;

  if (response === 0) {
    // User clicked "Update"
    await downloadUpdate(update.zipUrl);
    return true;
  } else if (response === 1) {
    // User clicked "Release Notes"
    shell.openExternal(`https://github.com/${REPO}/releases/latest`);
  }

  // For response === 2 (Later) or dialog closed (X button)
  return false;
}
