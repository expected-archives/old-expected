import React from 'react'
import ReactDOM from 'react-dom'
import createBrowserHistory from 'history/createBrowserHistory'
import { Avatar, Layout, Menu } from 'antd'
import { Link, Route, Router, Switch, Redirect } from 'react-router-dom'
import { ListContainer, ListImage, NewContainer } from './routes'
import './index.less'

const { Header, Content, Footer } = Layout

const App = () => (
    <Router history={createBrowserHistory()}>
        <Layout>
            <Header>
                <Menu theme={'dark'} mode={'horizontal'} style={{ lineHeight: '64px' }}>
                    <Menu.Item><Link to={'/containers'}>Containers</Link></Menu.Item>
                    <Menu.Item><Link to={'/images'}>Images</Link></Menu.Item>
                    <Menu.Item><Link to={'/containers/new'}>Create</Link></Menu.Item>

                    <Menu.Item style={{ float: 'right' }}>
                        <Avatar src={'https://avatars2.githubusercontent.com/u/32649258?v=4'} />
                    </Menu.Item>
                </Menu>
            </Header>
            <Content style={{ padding: '0 50px' }}>
                <div className={'content'}>
                    <Switch>
                        <Route name={'new-container'} path={'/containers/new'} component={NewContainer} />
                        <Route name={'list-container'} path={'/containers'} component={ListContainer} />
                        <Route name={'list-image'} path={'/images'} component={ListImage} />
                        <Redirect from={'/'} to={'/containers'} />
                    </Switch>
                </div>
            </Content>
            <Footer style={{ textAlign: 'center' }}>
                Expected.sh © 2019
            </Footer>
        </Layout>
    </Router>
)

ReactDOM.render(<App />, document.getElementById('root'))
