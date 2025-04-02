import { execSync } from "child_process";
import { readFileSync } from "fs";
import { fileURLToPath } from "url";
import path from "path";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const envFile = path.join(__dirname, "../.env.electron");
const token = readFileSync(envFile, "utf8")
  .match(/GH_TOKEN=(.*)/)[1]
  .trim();

process.env.GH_TOKEN = token;
execSync("npx electron-builder --publish always", { stdio: "inherit" });
