import {AccelerometerData} from './AccelerometerData';
import {LocationData} from './LocationData';
import {RiskData} from './RiskData';
import {Event} from './Event';

export interface Source {
    id: string;
    device_id: string;
    event_id: string;
    event_type: string;
    serial_number: string;
    message: string;
    max_prob: number;
    organization_id: string;
    hidden: boolean;
    start_time: string;
    end_time: string;
    incident_time: string;
    created_at: string;
    updated_at: string;
    thumbnail_key: string;
    video_key: string;
    accelerometer: AccelerometerData[];
    location: LocationData[];
    risk: RiskData[];
    events: Event[];
}
