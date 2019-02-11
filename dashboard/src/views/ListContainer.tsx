import React, {Component} from "react";
import TimeAgo from "react-timeago";
import {Link} from "react-router-dom";
import {Header, TableCard} from "../components";
import {fetchContainers, IContainer} from "../client";

interface IProps {
}

interface IState {
    containers: IContainer[];
}

export default class ListContainer extends Component<IProps, IState> {
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
            key: 'createdAt',
            render: (createdAt: any) => (
                <TimeAgo date={createdAt} minPeriod={10}/>
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

    constructor(props: IProps) {
        super(props);

        this.state = {
            containers: [],
        };
    }

    componentDidMount = () => {
        fetchContainers()
            .then((containers) => this.setState({containers}))
            .catch(console.error);
    };

    render = () => (
        <div>
            <Header title={'Containers'} pretitle={'Overview'}>
                <Link to={'/containers/new'} className={'btn btn-primary'}>
                    Create
                </Link>
            </Header>

            <TableCard<IContainer> columns={ListContainer.columns} dataSource={this.state.containers}
                                   onRowClick={data => console.log(data)}/>
        </div>
    );
}
