import React from 'react'
import ReactDOM from 'react-dom'
import createBrowserHistory from 'history/createBrowserHistory'
import { Avatar, Layout, Menu } from 'antd'
import { Link, Route, Router, Switch, Redirect } from 'react-router-dom'
import { ListContainer, ListImage } from './routes'
import './index.less'

const { Header, Content, Footer } = Layout

const App = () => (
    <Router history={createBrowserHistory()}>
        <Layout>
            <Header>
                <Menu theme={'dark'} mode={'horizontal'} style={{ lineHeight: '64px' }}>
                    <Menu.Item><Link to={'/containers'}>Containers</Link></Menu.Item>
                    <Menu.Item><Link to={'/images'}>Images</Link></Menu.Item>

                    <Menu.Item style={{ float: 'right' }}>
                        <Avatar src={'https://avatars2.githubusercontent.com/u/32649258?v=4'} />
                    </Menu.Item>
                </Menu>
            </Header>
            <Content style={{ padding: '0 50px' }}>
                <div style={{
                    background: '#fff',
                    marginTop: '64px',
                    padding: 24,
                    minHeight: 280,
                }}>
                    <Switch>
                        <Route name={'list-containers'} path={'/containers'} component={ListContainer} />
                        <Route name={'list-images'} path={'/images'} component={ListImage} />
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
