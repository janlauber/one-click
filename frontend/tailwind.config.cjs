/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    "./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}"
  ],
  theme: {
    extend: {
      colors: {
        primary: //"#3E6ABB",
        {
          100: "#B5B5B5",
          200: "#B5B5B5",
          300: "#B5B5B5",
          400: "#B5B5B5",
          500: "#0e0e0e",
          600: "#0e0e0e",
          700: "#0e0e0e",
          800: "#0e0e0e",
          900: "#0e0e0e"
        },
        secondary: //"#EDF1FF",
        {
          100: "#FFFFFF",
          200: "#FFFFFF",
          300: "#FFFFFF",
          400: "#FFFFFF",
          500: "#EDF1FF",
          600: "#C4CFFF",
          700: "#9BA9FF",
          800: "#727EFF",
          900: "#4959FF"
        },
        secondarydark: "#122036",
        text: "#0A111F",
        background: "#F7F9FC"
      },
      screens: {
        "hover-hover": { raw: "(hover: hover)" }
      },
      height: {
        128: "32rem"
      },
      keyframes: {
        wiggle: {
          "0%, 100%": { transform: "rotate(-0.5deg)" },
          "50%": { transform: "rotate(0.5deg)" }
        }
      },
      animation: {
        wiggle: "wiggle 1s ease-in-out infinite"
      }
    }
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("tailwind-scrollbar"),
    require("flowbite/plugin")
  ]
};
