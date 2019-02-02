import React from 'react'

export default ({ title, children }) => (
    <div className={'card'}>
        {title && (
            <div className={'card-header'}>
                <h4>{title}</h4>
            </div>
        )}
        <div className={'card-table table-responsive'}>
            <table className={'table'}>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>URL</th>
                        <th>Created</th>
                        <th>Tags</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>Hello</td>
                        <td>hello.remicaumette.ctr.expected.sh</td>
                        <td>1 minute ago</td>
                        <td></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
)
