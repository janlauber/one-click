import { writable } from "svelte/store";
import { browser } from "$app/environment";

const defaultValue = "";
let selectedDeploymentIdValue = defaultValue;

// Access localStorage only if running in the browser
if (browser) {
    const storedValue = localStorage.getItem("selectedDeploymentId");
    selectedDeploymentIdValue = storedValue !== null ? storedValue : defaultValue;
}

const selectedDeploymentId = writable<string>(selectedDeploymentIdValue);

selectedDeploymentId.subscribe((value) => {
    if (browser) {
        localStorage.setItem("selectedDeploymentId", value);
    }
});

export default selectedDeploymentId;
