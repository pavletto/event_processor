export interface Event {
    id: string;

    event_id: string;
    start_time: number;
    end_time: number;
    tag: 'danger' | 'safe';
    comment: string;
    created_at: string;
    updated_at: string;
}