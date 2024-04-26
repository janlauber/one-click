import { BookDashed, Boxes, Folder, Home } from "lucide-svelte";
import { createWritableStore } from "./data";
import { client } from "$lib/pocketbase";

export const breadcrumbItems = createWritableStore<any[]>([
    {
        title: "Home",
        href: "/app",
        icon: Home
    }
]);

export async function generateBreadcrumb(
    pathname: string
): Promise<Array<{ title: string; href: string; icon: any }>> {
    const segments = pathname.split("/").filter((segment) => segment !== "");
    let accumulatedPath = "";
    const breadcrumbItems = [];
    for (let i = 0; i < segments.length; i++) {
        accumulatedPath += `/${segments[i]}`;
        // Skip 'projects', 'deployments', and 'blueprints'
        if (["projects", "deployments", "blueprints"].includes(segments[i])) {
            continue;
        }
        let icon;
        if (breadcrumbItems.length === 0) {
            icon = Home; // Set the icon to Home for the first item
        } else {
            icon = determineIcon(segments[i - 1]); // Use the previous segment to determine the icon
        }
        let title = await getNameById(segments[i - 1], segments[i]); // Fetch the name based on the ID
        if (!title) {
            title = segments[i];
        }
        breadcrumbItems.push({
            title: title,
            href: accumulatedPath,
            icon: icon
        });
    }
    return breadcrumbItems;
}

function determineIcon(segment: string): any {
    switch (segment) {
        case "projects":
            return Folder;
        case "deployments":
            return Boxes;
        case "blueprints":
            return BookDashed;
        case "app":
            return Home;
        default:
            return undefined;
    }
}

async function getNameById(segment: string, id: string): Promise<string | undefined> {
    if (segment === "projects") {
        const response = await client.collection("projects").getOne(id);
        return response?.name;
    }
    if (segment === "deployments") {
        const response = await client.collection("deployments").getOne(id);
        return response?.name;
    }
    return undefined;
}
