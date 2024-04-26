// Interface for rollout events
export interface Event {
    reason: string;
    message: string;
    typus: string;
}

// Interface for rollout events response
export interface RolloutEventsResponse {
    events: Event[];
}
