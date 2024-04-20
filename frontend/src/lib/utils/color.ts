export const determineRolloutColor = (status: string) => {
    switch (status) {
        case "Pending":
            return "yellow";
        case "Not Ready":
            return "yellow";
        case "Error":
            return "red";
        case "OK":
            return "green";
        default:
            return "gray";
    }
};
