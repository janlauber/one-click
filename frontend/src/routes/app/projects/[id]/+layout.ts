import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import type { PageLoad } from "../../$types";

export const load: PageLoad = async ({ params, url }: any) => {
    const { id } = params;

    // use regex to match only exact or "/" paths
    const projectPathRegex = new RegExp(`/app/projects/${id}(/)?$`);

    // only update the projects if the current path matches the regex
    if (projectPathRegex.test(url.pathname)) {
        await updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: id
        }).catch(console.error);
    }
};
