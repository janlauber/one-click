import type { FrameworksResponse } from "$lib/pocketbase/generated-types";

export function frameworkLogoUrl(framework: FrameworksResponse | undefined): string {
    if (!framework) {
        return "";
    }
    return (
        "/api/files/" +
        framework?.collectionId +
        "/" +
        framework?.id +
        "/" +
        framework?.logo
    );
}
