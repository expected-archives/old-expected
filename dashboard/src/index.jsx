import React from 'react'
import ReactDOM from 'react-dom'
import createBrowserHistory from 'history/createBrowserHistory'
import { Link, Router } from 'react-router-dom'
import { Header, TableCard, Card } from './components'
import './styles/index.scss'

const App = () => (
    <Router history={createBrowserHistory()}>
        <div>
            <nav className={'navbar navbar-expand-lg navbar-dark'}>
                <div className={'container'}>
                    <Link to={'/'} className={'navbar-brand'}>Expected.sh</Link>
                    <button className={'navbar-toggler'} type={'button'}>
                        <span className={'navbar-toggler-icon'}/>
                    </button>
                    <div className={'collapse navbar-collapse'}>
                        <ul className={'navbar-nav'}>
                            <li className={'nav-item active'}>
                                <Link to={'/containers'} className={'nav-link'}>Containers</Link>
                            </li>
                            <li className={'nav-item'}>
                                <Link to={'/images'} className={'nav-link'}>Images</Link>
                            </li>
                        </ul>
                        <ul className={'navbar-nav ml-auto'}>
                            <li className={'nav-item'}>
                                <Link to={'/settings'} className={'nav-link'}>Settings</Link>
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

                <TableCard />
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
