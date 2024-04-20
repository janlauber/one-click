export type Color =
    | "none"
    | "red"
    | "yellow"
    | "green"
    | "indigo"
    | "purple"
    | "pink"
    | "blue"
    | "dark"
    | "primary"
    | undefined;

export function getTagColor(tag: string): Color {
    // available colors: "red" | "yellow" | "green" | "indigo" | "purple" | "blue" | "dark" | "none" | "pink" | "primary" | undefined
    // lookup if there are any of the following tokens in the tag: 'prod', 'dev', 'test', 'qa', 'staging', 'demo', 'int', 'uat'
    const color = tag.toLowerCase().match(/prod|dev|test|qa|staging|demo/)?.[0];

    const colorMap: Record<string, Color> = {
        prod: "red",
        dev: "yellow",
        test: "green",
        qa: "indigo",
        staging: "purple",
        demo: "blue",
        int: "pink",
        uat: "dark"
    };

    return color ? colorMap[color] : "none";
}
