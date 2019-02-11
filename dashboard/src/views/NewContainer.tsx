import React, {Component, FormEvent} from "react";
import {Form, FormGroup, Header} from "../components";
import {Link} from "react-router-dom";
import {createContainer} from "../client";

interface IProps {
}

interface IState {
    name: string;
    image: string;
    size: string;
    tags: string;
}

export default class NewContainer extends Component<IProps, IState> {
    constructor(props: IProps) {
        super(props);

        this.state = {
            name: "",
            image: "",
            size: "",
            tags: "",
        };
    }

    handleChange = (event: any) => {
        this.setState({[event.target.name]: event.target.value} as any)
    };

    handleSubmit = (event: FormEvent) => {
        event.preventDefault();

        const {name, image, size, tags } = this.state;
        createContainer(name, image, size, [])
            .then(console.log)
            .catch(console.error);
    };

    render = () => (
        <div className={'row justify-content-center'}>
            <div className={'col-12 col-lg-10 col-xl-8'}>
                <Header pretitle={'Containers'} title={'Create a new container'}/>

                <Form onSubmit={this.handleSubmit}>
                    <FormGroup name={'Name'}>
                        <input type="text" className="form-control"
                               placeholder={'my-container'} name={'name'} onChange={this.handleChange}/>
                    </FormGroup>

                    <FormGroup name={'Image'}>
                        <input type="text" className="form-control"
                               placeholder={'nginx:latest'} name={'image'} onChange={this.handleChange}/>
                    </FormGroup>

                    <FormGroup name={'Select a size'}>
                        <select className="form-control" name={'size'} onChange={this.handleChange}>
                            <option value={'64'}>64mb</option>
                            <option value={'128'}>128mb</option>
                            <option value={'256'}>256mb</option>
                        </select>
                    </FormGroup>

                    <FormGroup name={'Tags'}
                               description={'This is how others will learn about the project, so make it good!'}>
                        <input type="text" className="form-control" name={'tags'}
                               onChange={this.handleChange}/>
                    </FormGroup>

                    <hr className="mt-5 mb-5"/>

                    <button className={'btn btn-block btn-primary'}>
                        Create container
                    </button>
                    <Link to={'/containers'} className={'btn btn-block btn-link text-muted'}>
                        Cancel
                    </Link>
                </Form>
            </div>
        </div>
    );
}
