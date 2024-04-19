import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import type { PageLoad } from "../../$types";

export const load: PageLoad = async ({ params }: any) => {
    const { id } = params;

    // use regex to match both /projects/:id and /projects/:id/
    const projectPathRegex = new RegExp(`/app/projects/${id}/?$`, "i");

    // only update the projects if the current path matches the regex
    if (projectPathRegex.test(window.location.pathname)) {
        await updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: id
        }).catch(console.error);
    }
};
