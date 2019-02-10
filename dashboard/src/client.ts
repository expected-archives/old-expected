import axios from "axios";

const client = axios.create({
    baseURL: process.env.API_URL || "http://localhost:3000",
    headers: {
        Authorization: document.cookie.split('=')[1],
    },
});

const remapFields = (obj: any, fields: any) =>
    Object.entries(obj)
        .map(([key, value]) => [fields[key] || key, value ])
        .reduce((obj, [key, value]) => ({...obj, [key]: value }), {});

export interface IContainer {
    key: string;
    name: string;
    status: string;
    image: string;
    endpoint: string;
    memory: number;
    tags: string[];
    createdAt: Date;
}

export const fetchContainers = (): Promise<IContainer[]> =>
    client.get("/v1/containers")
        .then((res) => {
            if (res.status !== 200)
                throw new Error(res.data.message);
            return res.data.containers.map((container: object) =>
                remapFields(container, {created_at: "createdAt"}));
        });
