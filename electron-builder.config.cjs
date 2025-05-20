const isWindows = process.platform === "win32";
const isLinux = process.platform === "linux";

module.exports = {
  appId: "gg.quicksave.app",
  productName: "quicksave",
  asar: true,
  asarUnpack: ["assets"],
  directories: {
    output: "release/${version}",
  },
  files: ["dist-electron", "dist", "assets"],
  extraFiles: isWindows
    ? [
        {
          from: "backend/quicksaveService.exe",
          to: "../win-unpacked/backend/quicksaveService.exe",
        },
        {
          from: "backend/updater/updater.exe",
          to: "../win-unpacked/backend/updater.exe",
        },
      ]
    : isLinux
      ? [
          {
            from: "backend/quicksaveService",
            to: "../linux-unpacked/backend/quicksaveService",
          },
          {
            from: "backend/updater/updater",
            to: "../linux-unpacked/backend/updater",
          },
        ]
      : [],
  mac: {
    artifactName: "${productName}_${version}.${ext}",
    target: ["dmg", "zip"],
  },
  win: {
    artifactName: "${productName}.${ext}",
    icon: "assets/icon.ico",
    target: ["dir"],
  },
  linux: {
    target: ["AppImage"],
    icon: "assets/512x512.png",
    category: "Game",
    artifactName: "${productName}.${ext}",
  },
  publish: {
    provider: "github",
    owner: "Jehan1241",
    repo: "quicksave",
    releaseType: "release",
    publishAutoUpdate: true,
  },
};
