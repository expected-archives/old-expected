import React from 'react'

export default ({ title, children }) => (
    <div className={'card'}>
        {title && (
            <div className={'card-header'}>
                <h4>{title}</h4>
            </div>
        )}
        <div className={'card-body'}>
            {children}
        </div>
    </div>
)
