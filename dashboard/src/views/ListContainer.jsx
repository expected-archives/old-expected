import React from 'react'
import TimeAgo from 'react-timeago'
import { Link } from 'react-router-dom'
import { Header, TableCard } from '../components'

export default class ListContainer extends React.Component {
    static columns = [
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
            key: 'created_at',
            render: createdAt => (
                <TimeAgo date={createdAt} minPeriod={10} />
            ),
        },
        {
            title: 'Tags',
            key: 'tags',
            render: tags => (
                <span>
                    {tags.map((tag, index) => <div key={index}>{tag}</div>)}
                </span>
            ),
        },
    ]

    static dataSource = [
        {
            key: '1',
            name: 'Nginx',
            status: 'stopped',
            image: 'nginx:latest',
            endpoint: 'nginx.remicaumette.ctr.expected.sh',
            memory: 64,
            tags: ['frontend', 'expected.sh'],
            created_at: new Date(),
        },
        {
            key: '2',
            name: 'Mysql',
            status: 'started',
            image: 'mysql:latest',
            endpoint: 'mysql.remicaumette.ctr.expected.sh',
            memory: 128,
            tags: ['database', 'expected.sh'],
            created_at: new Date(Date.now() - 7200),
        },
    ]

    render() {
        return (
            <div>
                <Header title={'Containers'} pretitle={'Overview'}>
                    <Link to={'/containers/new'} className={'btn btn-primary'}>
                        Create
                    </Link>
                </Header>

                <TableCard dataSource={ListContainer.dataSource} columns={ListContainer.columns}
                            onRowClick={data => console.log(data)} />
            </div>
        )
    }
}
