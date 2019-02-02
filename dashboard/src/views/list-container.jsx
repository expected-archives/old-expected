import React from 'react'
// import TimeAgo from 'react-timeago'
import { Header, TableCard } from '../components'
import { Link } from 'react-router-dom'

// const dataSource = [{
//     key: '1',
//     name: 'Nginx',
//     status: 'stopped',
//     image: 'nginx:latest',
//     memory: 64,
//     tags: ['frontend', 'expected.sh'],
//     deploy_at: new Date(),
//     created_at: new Date(),
// }, {
//     key: '2',
//     name: 'Mysql',
//     status: 'started',
//     image: 'mysql:latest',
//     memory: 128,
//     tags: ['database', 'expected.sh'],
//     deploy_at: new Date(),
//     created_at: new Date(Date.now() - 7200),
// }]

// const columns = [{
//     title: 'Name',
//     dataIndex: 'name',
//     key: 'name',
// }, {
//     title: 'Status',
//     dataIndex: 'status',
//     key: 'status',
//     render: tag => (
//         <Tag
//             color={tag === 'started' ? 'green' : 'red'}>{tag === 'started' ? 'Started' : 'Stopped'}</Tag>
//     ),
// }, {
//     title: 'Image',
//     dataIndex: 'image',
//     key: 'image',
// }, {
//     title: 'Memory',
//     dataIndex: 'memory',
//     key: 'memory',
//     render: memory => `${memory}mb`,
// }, {
//     title: 'Last deployment',
//     dataIndex: 'deploy_at',
//     key: 'deploy_at',
//     render: deployAt => (
//         <TimeAgo date={deployAt} minPeriod={10}/>
//     ),
// }, {
//     title: 'Tags',
//     dataIndex: 'tags',
//     key: 'tags',
//     render: tags => (
//         <span>
//             {tags.map(tag => <Tag color={'blue'} key={tag}>{tag}</Tag>)}
//         </span>
//     ),
// }]

export default () => (
    <div>
        <Header title={'Containers'} pretitle={'Overview'}>
            <Link to={'/containers/new'} className={'btn btn-primary'}>
                Create
            </Link>
        </Header>

        <TableCard/>
    </div>
)
