import React from 'react'
import ReactDOM from 'react-dom'
import { Avatar, Layout, Menu, Breadcrumb } from 'antd'
import './index.less'

const { Header, Content, Footer } = Layout

const App = () => (
    <Layout>
        <Header>
            <Menu theme={'dark'} mode={'horizontal'} style={{ lineHeight: '64px' }}>
                <Menu.Item>Containers</Menu.Item>
                <Menu.Item>Images</Menu.Item>

                <Menu.Item style={{ float: 'right' }}>
                    <Avatar src={'https://avatars2.githubusercontent.com/u/32649258?v=4'} />
                </Menu.Item>
                <Menu.Item style={{ float: 'right' }}>Create</Menu.Item>
            </Menu>
        </Header>
        <Content style={{ padding: '0 50px' }}>
            <Breadcrumb style={{ margin: '16px 0' }}>
                <Breadcrumb.Item>Home</Breadcrumb.Item>
                <Breadcrumb.Item>Containers</Breadcrumb.Item>
            </Breadcrumb>
            <div style={{
                background: '#fff',
                padding: 24,
                minHeight: 280,
            }} />
        </Content>
        <Footer style={{ textAlign: 'center' }}>
            Expected.sh © 2019
        </Footer>
    </Layout>
)


ReactDOM.render(<App />, document.getElementById('root'))
