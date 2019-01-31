import React from 'react'
import { Form, Input, Select, Button } from 'antd'

const formItemLayout = {}

export default (props) => {
    const onSubmit = () => {
        console.log(props)
    }

    return (
        <div>
            <h1>Create Containers</h1>
            <Form onSubmit={onSubmit}>
                <Form.Item {...formItemLayout} label={'Name'} required>
                    <Input type={'text'} maxLength={30} />
                </Form.Item>
                <Form.Item {...formItemLayout} label={'Image'} required>
                    <Input type={'text'} maxLength={255} />
                </Form.Item>
                <Form.Item {...formItemLayout} label={'Size'} required>
                    <Select defaultValue={64}>
                        <Select.Option value={64}>64mb</Select.Option>
                        <Select.Option value={128}>128mb</Select.Option>
                        <Select.Option value={256}>256mb</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item {...formItemLayout}>
                    <Button type={'primary'} htmlType={'submit'}>Create</Button>
                </Form.Item>
            </Form>
        </div>
    )
}
