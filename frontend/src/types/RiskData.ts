export interface RiskData {
    id: number;
    event_id: string;
    recording_time: string;
    prob: number;

    relative_time: Date;
    recording_time_date: Date;
}
