export const initTileSize = (setTileSize: (size: number) => void) => {
  const tileSize = Number(localStorage.getItem("tileSize"));
  if (tileSize !== 0) {
    setTileSize(tileSize);
  } else {
    setTileSize(35);
    localStorage.setItem("tileSize", "35");
  }
};

export const setLastPath = (navigate: any) => {
  const lastPath = localStorage.getItem("lastPath");
  if (!lastPath) {
    navigate("/library");
  } else {
    navigate(lastPath);
  }
};
