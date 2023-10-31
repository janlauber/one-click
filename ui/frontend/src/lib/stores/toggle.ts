import { writable } from "svelte/store";

const defaultValue = false;
const showSurveyValue =
    localStorage.getItem("showSurvey") === null
        ? defaultValue
        : localStorage.getItem("showSurvey") === "true";

const showSurvey = writable<boolean>(showSurveyValue);

showSurvey.subscribe((value) => {
    localStorage.setItem("showSurvey", value.toString());
});

export default showSurvey;
