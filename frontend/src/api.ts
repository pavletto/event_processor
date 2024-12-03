import axios from 'axios';
import {Source} from "./types/Source.ts";

export const api = axios.create({
    baseURL: import.meta.env.PROD ? 'http://localhost:8080/' : '/api',
})

export const sources = () => {
    return api.get<Source[]>('/sources')
}
export const source = (id: string) => {
    return api.get<Source>(`/source/${id}`)
}
export const upload = (formData: FormData) => {
    axios
        .post('http://localhost:8080/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        })
}