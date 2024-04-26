import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
    await updateDataStores({
        filter: UpdateFilterEnum.ALL
    }).catch((error) => {
        console.error(error);
    });
};
