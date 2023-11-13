import type { TechnologiesResponse } from "$lib/pocketbase/generated-types";

export function technologyLogoUrl(technology: TechnologiesResponse | undefined): string {
    if (!technology) {
        return "";
    }
    return (
        "/api/files/" +
        technology?.collectionId +
        "/" +
        technology?.id +
        "/" +
        technology?.logo
    );
}
