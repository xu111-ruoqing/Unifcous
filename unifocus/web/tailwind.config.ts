import type { Config } from "tailwindcss";
import tailwindcssAnimate from "tailwindcss-animate";

const config: Config = {
  darkMode: ["class"],
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // shadcn CSS variable mappings
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
          border: "#E4DDD0",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        // Warm academic palette
        page: "#F3EEE5",
        sidebar: { bg: "#1F1B17", accent: "#C4622D" },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
          indigo: "#C4622D",   // sienna — primary active
          cyan:   "#3A7555",   // sage green
          amber:  "#C8880A",   // warm amber
          violet: "#3E5672",   // slate blue
          rose:   "#B84040",
        },
        // Kept for planet view (uses inline styles anyway)
        banner: { from: "#12103a", via: "#2d1b69", to: "#6c5ce7" },
        text: {
          primary:   "#1A1714",
          secondary: "#6A5E53",
          muted:     "#A89B8A",
        },
      },
      fontFamily: {
        display: ["Cormorant Garamond", "Songti SC", "Georgia", "serif"],
        body:    ["Sora", "PingFang SC", "sans-serif"],
        mono:    ["IBM Plex Mono", "monospace"],
      },
      borderRadius: { sm: "6px", md: "10px", lg: "16px" },
    },
  },
  plugins: [tailwindcssAnimate],
};

export default config;
