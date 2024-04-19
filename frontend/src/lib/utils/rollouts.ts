import type { RolloutMetricsResponse } from "$lib/types/metrics";
import type { RolloutStatusResponse } from "$lib/types/status";

async function fetchFromAPI(endpoint: string) {
    let data: RolloutStatusResponse | RolloutMetricsResponse | RolloutEventsResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return data;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    const baseUrl = window.location.hostname === "localhost" ? `http://localhost:8090` : '';
    const url = `${baseUrl}/pb/${endpoint}`;

    try {
        const response = await fetch(url, {
            headers: authHeader
        });
        data = await response.json();
    } catch (error) {
        // Handle error
    }

    return data;
}

export async function getRolloutStatus(projectId: string, deploymentId: string) {
    return fetchFromAPI(`${projectId}/${deploymentId}/status`);
}

export async function getRolloutMetrics(projectId: string, deploymentId: string) {
    return fetchFromAPI(`${projectId}/${deploymentId}/metrics`);
}

export async function getRolloutEvents(projectId: string, deploymentId: string) {
    return fetchFromAPI(`${projectId}/${deploymentId}/events`);
}
