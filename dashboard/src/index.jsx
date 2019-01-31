import React from 'react'
import ReactDOM from 'react-dom'
import createBrowserHistory from 'history/createBrowserHistory'
import { Link, Router } from 'react-router-dom'
import { Header, Card } from './components'
import './index.scss'

const App = () => (
    <Router history={createBrowserHistory()}>
        <div>
            <nav className="navbar navbar-expand-lg navbar-dark">
                <div className={'container'}>
                    <Link to={'/'} className={'navbar-brand'}>Expected.sh</Link>
                    <button className="navbar-toggler" type="button" data-toggle="collapse"
                            data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false"
                            aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarNav">
                        <ul className="navbar-nav">
                            <li className="nav-item active">
                                <a className="nav-link" href="#">Containers</a>
                            </li>
                            <li className="nav-item">
                                <a className="nav-link" href="#">Images</a>
                            </li>
                        </ul>
                        <ul className="navbar-nav ml-auto">
                            <li className="nav-item">
                                <a className="nav-link" href="#">Settings</a>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>

            <div className={'container'}>
                <Header title={'Containers'} pretitle={'Overview'}>
                    <Link to={'/containers/new'} className={'btn btn-primary'}>
                        Create
                    </Link>
                </Header>

                <Card title={'Hello'}>
                    <p>Hello</p>
                </Card>
            </div>
        </div>
    </Router>
)

ReactDOM.render(<App/>, document.getElementById('root'))

// <Layout>
// <Header>
// <Menu theme={'dark'} mode={'horizontal'} style={{ lineHeight: '64px' }}>
// <Menu.Item><Link to={'/containers'}>Containers</Link></Menu.Item>
// <Menu.Item><Link to={'/images'}>Images</Link></Menu.Item>
// <Menu.Item><Link to={'/containers/new'}>Create</Link></Menu.Item>
//
// <Menu.Item style={{ float: 'right' }}>
// <Avatar src={'https://avatars2.githubusercontent.com/u/32649258?v=4'} />
// </Menu.Item>
// </Menu>
// </Header>
// <Content style={{ padding: '0 50px' }}>
// <div className={'content'}>
//     <Switch>
//     <Route name={'new-container'} path={'/containers/new'} component={NewContainer} />
// <Route name={'list-container'} path={'/containers'} component={ListContainer} />
// <Route name={'list-image'} path={'/images'} component={ListImage} />
// <Redirect from={'/'} to={'/containers'} />
// </Switch>
// </div>
// </Content>
// <Footer style={{ textAlign: 'center' }}>
//     Expected.sh © 2019
// </Footer>
// </Layout>
