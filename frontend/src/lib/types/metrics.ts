// Interface for individual pod metrics
interface PodMetrics {
    name: string;
    cpu: string;    // Assuming the CPU usage is represented as a string
    memory: string; // Memory usage as a string (e.g., "8444Ki")
}

// Interface for the overall metrics response
export interface RolloutMetricsResponse {
    metrics: PodMetrics[];
}
