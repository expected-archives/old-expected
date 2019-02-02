import React from 'react'
import { Header } from '../components'
import { Link } from 'react-router-dom'

const formItemLayout = {}

export default (props) => {
    const onSubmit = () => {
        console.log(props)
    }

    return (
        <div className={'row justify-content-center'}>
            <div className={'col-12 col-lg-10 col-xl-8'}>
                <Header pretitle={'Containers'} title={'Create a new container'}/>

                <form className={'mb-5'}>
                    <div className="form-group">
                        <label>Name</label>
                        <input type="text" className="form-control" placeholder={'my-container'}/>
                    </div>

                    <div className="form-group">
                        <label>Image</label>
                        <input type="text" className="form-control" placeholder={'nginx:latest'}/>
                    </div>

                    <div className="form-group">
                        <label>Select a size</label>
                        <select className="form-control">
                            <option value={'64'}>64mb</option>
                            <option value={'128'}>128mb</option>
                            <option value={'256'}>256mb</option>
                        </select>
                    </div>

                    <div className="form-group">
                        <label>Tags</label>
                        <input type="text" className="form-control"/>
                    </div>

                    <hr className="mt-5 mb-5" />

                    <button className={'btn btn-block btn-primary'}>
                        Create container
                    </button>
                    <Link to={'/containers'} className={'btn btn-block btn-link text-muted'}>
                        Cancel
                    </Link>
                </form>
            </div>
        </div>
    )
}
