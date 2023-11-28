// Interface for rollout events
interface Event {
    Reason: string;
    Message: string;
    Type: string;
}

// Interface for rollout events response

interface RolloutEventsResponse {
    events: Event[];
}
