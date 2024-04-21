import { browser } from "$app/environment";
import { client } from "$lib/pocketbase";
import { breadcrumbItems, generateBreadcrumb } from "$lib/stores/breadcrumb";
import { redirect } from "@sveltejs/kit";

// turn off SSR - we're JAMstack here
export const ssr = false;
// Prerendering turned off. Turn it on if you know what you're doing.
export const prerender = false;
// trailing slashes make relative paths much easier
export const trailingSlash = "always";

export const load = async ({ url }) => {
    const { pathname } = url;

    console.log(url.pathname);

    // i want to dynamically change the breadcrumbItems store (array of title, href, and icon) based on the url.pathname
    // e.g. /app/projects/1/deployments/1 should have a breadcrumb of Home > Projects > Project 1 > Deployments > Deployment 1

    // if url.pathname is /app, breadcrumbItems should be [{title: "Home", href: "/app", icon: Home}]
    // or if url.pathname is /app/blueprints then it should have a breadcrumb of Home > Blueprints or something like that
    // but if the user goes from /app/projects/1 to /app/blueprints, the breadcrumb Home > Projects > Project 1 > Blueprints should be shown

    const breadcrumb = generateBreadcrumb(pathname);
    const updateBreadcrumb = async () => {
        await breadcrumbItems.set(await breadcrumb);
    };

    await updateBreadcrumb();

    if (browser) {
        if (client.authStore.model && client.authStore.isValid) {
            if (pathname === "/login/" || pathname === "/signup/") {
                redirect(307, "/app");
            }
        } else {
            // if pathname contains "/app" and user is not logged in, redirect to login
            if (pathname.includes("/app")) {
                redirect(307, "/login");
            }
        }
    }

    return {
        pathname
    };
};
