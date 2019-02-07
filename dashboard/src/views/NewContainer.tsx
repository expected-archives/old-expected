import React, { FormEvent, Component } from "react";
import { Form, FormGroup, Header } from '../components'
import { Link } from 'react-router-dom'


export default class NewContainer extends Component {
    onSubmit(event: FormEvent) {
        event.preventDefault()
        console.log(event)
    }

    render() {
        return (
            <div className={'row justify-content-center'}>
                <div className={'col-12 col-lg-10 col-xl-8'}>
                    <Header pretitle={'Containers'} title={'Create a new container'}/>

                    <Form onSubmit={this.onSubmit.bind(this)}>
                        <FormGroup name={'Name'}>
                            <input type="text" className="form-control"
                                   placeholder={'my-container'} ref={'name'}/>
                        </FormGroup>

                        <FormGroup name={'Image'}>
                            <input type="text" className="form-control"
                                   placeholder={'nginx:latest'} ref={'image'}/>
                        </FormGroup>

                        <FormGroup name={'Select a size'}>
                            <select className="form-control" ref={'size'}>
                                <option value={'64'}>64mb</option>
                                <option value={'128'}>128mb</option>
                                <option value={'256'}>256mb</option>
                            </select>
                        </FormGroup>

                        <FormGroup name={'Tags'}
                                   description={'This is how others will learn about the project, so make it good!'}>
                            <input type="text" className="form-control" ref={'tags'}/>
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
        )
    }
}
