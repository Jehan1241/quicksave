const { Sidebar } = require("lucide-react");

/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: [
    "./pages/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./src/**/*.{ts,tsx}",
  ],
  prefix: "",
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        sm: "640px",
        // => @media (min-width: 640px) { ... }

        md: "768px",
        // => @media (min-width: 768px) { ... }

        lg: "1024px",
        // => @media (min-width: 1024px) { ... }

        xl: "1280px",
        // => @media (min-width: 1280px) { ... }

        "2xl": "1400px",
        // => @media (min-width: 1536px) { ... }
      },
    },
    extend: {
      colors: {
        Sidebar: "hsl(var(--sidebar))",
        content: "hsl(var(--content))",
        sliderEmpty: "hsl(var(--sliderEmpty))",
        sliderFull: "hsl(var(--sliderFull))",
        sliderThumb: "hsl(var(--sliderThumb))",
        viewButtons: "hsl(var(--viewButtons))",
        activeViewButton: "hsl(var(--activeViewButton))",
        hoverViewButton: "hsl(var(--hoverViewButton))",
        topBarButtons: "hsl(var(--topBarButtons))",
        topBarButtonsHover: "hsl(var(--topBarButtonsHover))",
        selectedSortItem: "hsl(var(--selectedSortItem))",
        selectedSortItemText: "hsl(var(--selectedSortItemText))",
        selectedSortItemHover: "hsl(var(--selectedSortItemHover))",
        dialogSaveButtons: "hsl(var(--dialogSaveButtons))",
        dialogSaveButtonsHover: "hsl(var(--dialogSaveButtonsHover))",
        dialogSaveButtonsText: "hsl(var(--dialogSaveButtonsText))",
        leftbarIcons: "hsl(var(--leftbarIcons))",
        playButton: "hsl(var(--playButton))",
        playButtonHover: "hsl(var(--playButtonHover))",
        playButtonText: "hsl(var(--playButtonText))",
        editButton: "hsl(var(--editButton))",
        editButtonHover: "hsl(var(--editButtonHover))",
        editButtonText: "hsl(var(--editButtonText))",
        platformBadge: "hsl(var(--platformBadge))",
        platformBadgeHover: "hsl(var(--platformBadgeHover))",
        platformBadgeText: "hsl(var(--platformBadgeText))",
        tagsBadge: "hsl(var(--tagsBadge))",
        tagsBadgeHover: "hsl(var(--tagsBadgeHover))",
        tagsBadgeText: "hsl(var(--tagsBadgeText))",
        devsBadge: "hsl(var(--devsBadge))",
        devsBadgeHover: "hsl(var(--devsBadgeHover))",
        devsBadgeText: "hsl(var(--devsBadgeText))",
        emptyGameTile: "hsl(var(--emptyGameTile))",
        emptyGameTileText: "hsl(var(--emptyGameTileText))",
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      fontFamily: {
        sans: ["Geist", "sans-serif"],
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [
    require("tailwindcss-animate"),
    require("tailwind-scrollbar")({
      nocompatible: true, // Enables non-standard utilities for more customization
      preferredStrategy: "pseudoelements", // Uses the pseudoelement strategy for scrollbars
    }),
  ],
};
