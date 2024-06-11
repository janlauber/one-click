// src/stores/version.js
import { writable } from "svelte/store";

export const version = writable(import.meta.env.VITE_APP_VERSION);
