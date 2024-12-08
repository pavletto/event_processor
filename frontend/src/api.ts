import axios from 'axios';
import {Source} from "@/types/Source";

export const api = axios.create({
    baseURL: import.meta.env.PROD && import.meta.env.VITE_HOST ? import.meta.env.VITE_HOST : '/api',
})

export const sources = () => {
    return api.get<Source[]>('/sources')
}
export const source = (id: string) => {
    return api.get<Source>(`/source/${id}`)
}
export const upload = (formData: FormData) => {
    return api.post<FormData>('/upload', formData, {headers: {
                'Content-Type': 'multipart/form-data',
            },
        })
}

export const createEvent = (event: Event) => {
    return axios.post('/api/events', event).then(() => {
        sources()
    })
}