import axios from "axios";

const client = axios.create({
    baseURL: process.env.API_URL || "http://localhost:3000",
    headers: {
        Authorization: document.cookie.split('=')[1],
    },
});

const remapFields = (obj: any, fields: any): any =>
    Object.entries(obj)
        .map(([key, value]) => [fields[key] || key, value])
        .reduce((obj, [key, value]) => ({...obj, [key]: value}), {});

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

const toContainer = (data: object): IContainer => {
    const container = remapFields(data, {created_at: "createdAt"});
    container.createdAt = new Date(container.createdAt);
    return container;
};

export const fetchContainers = (): Promise<IContainer[]> =>
    client.get("/v1/containers")
        .then((res) => {
            if (res.status !== 200)
                throw new Error(res.data.message);
            return res.data.containers.map((data: object) => toContainer(data));
        });

export const createContainer = (name: string, image: string, size: string, tags: string[]): Promise<IContainer> =>
    client.post("/v1/containers", { name, image, size, tags })
        .then((res) => {
            if (res.status !== 200)
                throw new Error(res.data.message);
            return toContainer(res.data.container);
        });
