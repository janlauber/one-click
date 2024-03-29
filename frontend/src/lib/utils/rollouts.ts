import type { RolloutMetricsResponse } from "$lib/types/metrics";
import type { RolloutStatusResponse } from "$lib/types/status";

export async function getRolloutStatus(projectId: string, rolloutId: string) {
    const response = await fetchRolloutStatus(projectId, rolloutId);
    return response;
}

// fetch from /rollouts/{projectId}/{rolloutId}/status
async function fetchRolloutStatus(projectId: string, rolloutId: string) {
    let status: RolloutStatusResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return status;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
        try {
            const response = await fetch(
                `http://localhost:8090/rollouts/${projectId}/${rolloutId}/status`,
                {
                    headers: authHeader
                }
            );
            status = await response.json();
        } catch (error) {
        }

        return status;
    }

    try {
        const response = await fetch(`/rollouts/${projectId}/${rolloutId}/status`, {
            headers: authHeader
        });
        status = await response.json();
    } catch (error) {
        // Handle error
    }

    return status;
}

export async function getRolloutMetrics(projectId: string, rolloutId: string) {
    const response = await fetchRolloutMetrics(projectId, rolloutId);
    return response;
}

async function fetchRolloutMetrics(projectId: string, rolloutId: string) {
    let metrics: RolloutMetricsResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return metrics;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
        try {
            const response = await fetch(
                `http://localhost:8090/rollouts/${projectId}/${rolloutId}/metrics`,
                {
                    headers: authHeader
                }
            );
            metrics = await response.json();
        } catch (error) {
            // Handle error silently
        }

        return metrics;
    }

    try {
        const response = await fetch(`/rollouts/${projectId}/${rolloutId}/metrics`, {
            headers: authHeader
        });
        metrics = await response.json();
    } catch (error) {
        // Handle error silently
    }

    return metrics;
}

// get rollout events
export async function getRolloutEvents(projectId: string, rolloutId: string) {
    const response = await fetchRolloutEvents(projectId, rolloutId);
    return response;
}

async function fetchRolloutEvents(projectId: string, rolloutId: string) {
    let events: RolloutEventsResponse | undefined;

    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
        return events;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };
    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
        try {
            const response = await fetch(
                `http://localhost:8090/rollouts/${projectId}/${rolloutId}/events`,
                {
                    headers: authHeader
                }
            );
            events = await response.json();
        } catch (error) {
            // Handle error
        }

        return events;
    }

    try {
        const response = await fetch(`/rollouts/${projectId}/${rolloutId}/events`, {
            headers: authHeader
        });
        events = await response.json();
    } catch (error) {
        // Handle error
    }

    return events;
}
