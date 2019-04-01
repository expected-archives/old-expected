import React from "react";
import TimeAgo from "react-timeago";
import { Link } from "react-router-dom";
import { Header, TableCard } from "../components";
import axios from "axios";

interface IContainerDataSource {
    key: string;
    name: string;
    status: string;
    image: string;
    endpoint: string;
    memory: number;
    tags: string[];
    createdAt: Date;
}

export default () => {
    const columns = [
        {
            title: 'Name',
            key: 'name',
        },
        {
            title: 'Endpoint',
            key: 'endpoint',
        },
        {
            title: 'Created',
            key: 'createdAt',
            render: (createdAt: any) => (
                <TimeAgo date={createdAt} minPeriod={10} />
            ),
        },
        {
            title: 'Tags',
            key: 'tags',
            render: (tags: any) => (
                <span>
                    {tags.map((tag: string, index: number) => (
                        <div key={index}>
                            {tag}
                        </div>
                    ))}
                </span>
            ),
        },
    ];

    const client = axios.create({
        baseURL: "http://localhost:3000",
        headers: {
            Authorization: document.cookie.split('=')[1],
        },
    })

    client.get("/v1/containers").then(console.log)

    const dataSource = [
        {
            key: '1',
            name: 'Nginx',
            status: 'stopped',
            image: 'nginx:latest',
            endpoint: 'nginx.remicaumette.ctr.expected.sh',
            memory: 64,
            tags: ['frontend', 'expected.sh'],
            createdAt: new Date(),
        },
        {
            key: '2',
            name: 'Mysql',
            status: 'started',
            image: 'mysql:latest',
            endpoint: 'mysql.remicaumette.ctr.expected.sh',
            memory: 128,
            tags: ['database', 'expected.sh'],
            createdAt: new Date(Date.now() - 7200),
        },
    ];

    return (
        <div>
            <Header title={'Containers'} pretitle={'Overview'}>
                <Link to={'/containers/new'} className={'btn btn-primary'}>
                    Create
                </Link>
            </Header>

            <TableCard<IContainerDataSource> dataSource={dataSource} columns={columns}
                onRowClick={data => console.log(data)} />
        </div>
    );
};
