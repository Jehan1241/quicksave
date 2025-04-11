const THEME_KEY = "theme";

export function setTheme() {
  const theme = localStorage.getItem(THEME_KEY);
  removeAllThemes();
  if (theme) {
    document.documentElement.classList.add(theme);
    localStorage.setItem(THEME_KEY, theme);
  } else {
    document.documentElement.classList.add("dark");
    localStorage.setItem(THEME_KEY, "dark");
  }
  console.log(theme);
}

export function updateTheme(theme: string) {
  switch (theme) {
    case "light":
      lightMode();
      break;
    case "red":
      redMode();
      break;
    case "magenta-dark":
      darkPurpleMode();
      break;
    case "dark":
      darkMode();
      break;
    default:
      darkMode();
      break;
  }
}

export function lightMode() {
  removeAllThemes();
  localStorage.setItem(THEME_KEY, "light");
}

export function darkPurpleMode() {
  removeAllThemes();
  document.documentElement.classList.add("magenta-dark");
  localStorage.setItem(THEME_KEY, "magenta-dark");
}

export function darkMode() {
  removeAllThemes();
  document.documentElement.classList.add("dark");
  localStorage.setItem(THEME_KEY, "dark");
}

export function redMode() {
  removeAllThemes();
  document.documentElement.classList.add("red");
  localStorage.setItem(THEME_KEY, "red");
}

function removeAllThemes() {
  document.documentElement.classList.remove("red");
  document.documentElement.classList.remove("dark");
  document.documentElement.classList.remove("magenta-dark");
}
