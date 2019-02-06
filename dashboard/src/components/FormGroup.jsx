import React from 'react'
import PropTypes from 'prop-types'

const FormGroup = ({ name, description, children }) => (
    <div className={'form-group'}>
        {name && (
            <label>{name}</label>
        )}
        {description && (
            <small className={'form-text text-muted'}>
                {description}
            </small>
        )}
        {children}
    </div>
)

FormGroup.propTypes = {
    name: PropTypes.string,
    description: PropTypes.string,
    children: PropTypes.node,
}

export default FormGroup
