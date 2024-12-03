export interface LocationData {
    id: number;
    event_id: string;
    recording_time: string;
    altitude?: number;
    latitude?: number;
    longitude?: number;
    speed_kmh?: number;
    speed_knots?: number;
    speed_mph?: number;
}
