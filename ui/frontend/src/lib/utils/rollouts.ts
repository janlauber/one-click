import type { RolloutStatusResponse } from "$lib/types/status";

export async function getRolloutStatus(projectId: string, rolloutId: string) {
    const response = await fetchRolloutStatus(projectId, rolloutId);
    return response;
}

// fetch from /rollouts/{projectId}/{rolloutId}/status
async function fetchRolloutStatus(projectId: string, rolloutId: string) {
    let status: RolloutStatusResponse | undefined;

    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
        try {
            const response = await fetch(
                `http://localhost:8090/rollouts/${projectId}/${rolloutId}/status`
            );
            status = await response.json();
        } catch (error) {
            console.error(error);
        }

        return status;
    }

    try {
        const response = await fetch(`/rollouts/${projectId}/${rolloutId}/status`);
        status = await response.json();
    } catch (error) {
        console.error(error);
    }

    return status;
}
