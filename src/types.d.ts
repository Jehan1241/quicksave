// This allows TypeScript to pick up the magic constants that's auto-generated by Forge's Vite
// plugin that tells the Electron app where to look for the Vite-bundled app code (depending on
// whether you're running in development or production).
declare interface Window {
  windowFunctions: {
    nukeCache: () => void;
    closeApp: () => void;
    minimize: () => void;
    maximize: () => void;
  };
  electron: {
    validateGamePath: (gamePath) => any;
    browseFileHandler: (
      options: Electron.OpenDialogOptions
    ) => Promise<Electron.OpenDialogReturnValue>;
  };
  themeMode: {
    toggle: () => Promise<boolean>;
    dark: () => Promise<void>;
    light: () => Promise<void>;
    system: () => Promise<boolean>;
    current: () => Promise<"dark" | "light" | "system">;
  };
}
