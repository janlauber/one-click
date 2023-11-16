import { writable } from "svelte/store";
import { browser } from "$app/environment";

const defaultValue = "";
let selectedProjectIdValue = defaultValue;

// Access localStorage only if running in the browser
if (browser) {
    const storedValue = localStorage.getItem("selectedProjectId");
    selectedProjectIdValue = storedValue !== null ? storedValue : defaultValue;
}

const selectedProjectId = writable<string>(selectedProjectIdValue);

selectedProjectId.subscribe((value) => {
    if (browser) {
        localStorage.setItem("selectedProjectId", value);
    }
});

export default selectedProjectId;
