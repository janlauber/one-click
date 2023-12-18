// Interface for rollout events
interface Event {
    reason: string;
    message: string;
    typus: string;
}

// Interface for rollout events response

interface RolloutEventsResponse {
    events: Event[];
}
