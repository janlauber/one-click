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
        primary: "#3E6ABB",
        secondary: "#EDF1FF",
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
    require("flowbite/plugin"),
    require("daisyui")
  ],
  daisyui: {
    logs: false,
    themes: ["emerald"]
  }
};
