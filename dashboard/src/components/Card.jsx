import React from 'react'
import PropTypes from 'prop-types'

const Card = ({ title, children }) => (
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

Card.propTypes = {
    title: PropTypes.string,
}

export default Card
