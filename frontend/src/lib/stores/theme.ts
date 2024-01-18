// get localStorage theme
import { writable } from "svelte/store";

const defaultValue = 'light';
const colorThemeValue: string | null =
    localStorage.getItem("color-theme");

const colorTheme = writable<string>(colorThemeValue || defaultValue);

colorTheme.subscribe((value) => {
    localStorage.setItem("color-theme", value.toString());
});

export default colorTheme;
