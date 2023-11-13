export function formatDateTime(dateTimeStr: string | undefined): string {
    if (!dateTimeStr) {
        return "";
    }

    // Create a new Date object using the input string
    const date = new Date(dateTimeStr);

    // Format the date and time components
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, "0");
    const day = date.getDate().toString().padStart(2, "0");
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    const seconds = date.getSeconds().toString().padStart(2, "0");

    // Return the formatted date-time string
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

export function timeAgo(dateTimeStr: string | undefined): string {
    if (!dateTimeStr) {
        return "";
    }

    // Create a new Date object using the input string
    const date = new Date(dateTimeStr);

    // Get the current date
    const now = new Date();

    // Get the difference in seconds between the two dates
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000);

    // Calculate nearest time unit (seconds, minutes, hours, days, weeks, months, years)
    const unit = diff < 60 ? "second" : diff < 3600 ? "minute" : diff < 86400 ? "hour" : diff < 604800 ? "day" : diff < 2629800 ? "week" : diff < 31557600 ? "month" : "year";

    // Calculate the number of units
    const num = Math.floor(diff / (unit === "second" ? 1 : unit === "minute" ? 60 : unit === "hour" ? 3600 : unit === "day" ? 86400 : unit === "week" ? 604800 : unit === "month" ? 2629800 : 31557600));

    // Return the formatted date-time string
    return `${num} ${unit}${num !== 1 ? "s" : ""} ago`;
}
