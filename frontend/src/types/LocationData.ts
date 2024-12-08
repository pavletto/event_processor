export interface LocationData {
    id: number;
    event_id: string;
    recording_time: string;
    relative_time: number;
    recording_time_date: Date;
    altitude?: number;
    latitude?: number;
    longitude?: number;
    speed_kmh?: number;
    speed_knots?: number;
    speed_mph?: number;
}
