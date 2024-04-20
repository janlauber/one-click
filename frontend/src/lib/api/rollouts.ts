import type { RolloutEventsResponse } from "$lib/types/events";
import type { RolloutMetricsResponse } from "$lib/types/metrics";
import type { RolloutStatusResponse } from "$lib/types/status";

async function fetchFromAPI<T>(endpoint: string): Promise<T | undefined> {
    let data: T | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return data;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    const baseUrl = window.location.hostname === "localhost" ? `http://localhost:8090` : "";
    const url = `${baseUrl}/pb/${endpoint}`;

    try {
        const response = await fetch(url, {
            headers: authHeader
        });
        data = (await response.json()) as T;
    } catch (error) {
        // Handle error
    }

    return data;
}

export async function getRolloutStatus(
    projectId: string,
    deploymentId: string
): Promise<RolloutStatusResponse | undefined> {
    return fetchFromAPI<RolloutStatusResponse>(`${projectId}/${deploymentId}/status`);
}

export async function getRolloutMetrics(
    projectId: string,
    deploymentId: string
): Promise<RolloutMetricsResponse | undefined> {
    return fetchFromAPI<RolloutMetricsResponse>(`${projectId}/${deploymentId}/metrics`);
}

export async function getRolloutEvents(
    projectId: string,
    deploymentId: string
): Promise<RolloutEventsResponse | undefined> {
    return fetchFromAPI<RolloutEventsResponse>(`${projectId}/${deploymentId}/events`);
}
