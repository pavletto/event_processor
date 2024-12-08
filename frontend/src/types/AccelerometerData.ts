export interface AccelerometerData {
    id: number;
    event_id: string;
    recording_time: string;
    relative_time: string;
    x: number;
    y: number;
    z: number;

    recording_time_date: Date;
    relative_time_number: number;
}
