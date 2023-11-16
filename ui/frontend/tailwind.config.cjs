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
          100: "#E6E9F4",
          200: "#C2CBE0",
          300: "#9EA9CC",
          400: "#7A86B8",
          500: "#5664A4",
          600: "#3E6ABB",
          700: "#324F91",
          800: "#253666",
          900: "#18213C"
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
    require("flowbite/plugin"),
    require("daisyui")
  ],
  daisyui: {
    logs: false,
    themes: ["emerald"]
  }
};
