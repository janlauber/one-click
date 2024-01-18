export function extractPaths(obj: any, prefix: string = ''): string[] {
    let paths: string[] = [];

    for (const key in obj) {
        if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
            // If the value is an object, recursively find paths
            paths = paths.concat(extractPaths(obj[key], `${prefix}${key}.`));
        } else {
            // If the value is not an object (or is an array), add the path to the list
            paths.push(`${prefix}${key}`);
        }
    }

    return paths;
}
