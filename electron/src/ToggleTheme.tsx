const THEME_KEY = "theme";

export function lightMode() {
    document.documentElement.classList.remove("red");
    document.documentElement.classList.remove("dark");
    localStorage.setItem(THEME_KEY, "light");
}

export function darkMode() {
    document.documentElement.classList.add("dark");
    document.documentElement.classList.remove("red");
    localStorage.setItem(THEME_KEY, "dark");
}

export function redMode() {
    document.documentElement.classList.add("red");
    document.documentElement.classList.remove("dark");
    localStorage.setItem(THEME_KEY, "light");
}
