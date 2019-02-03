import React from 'react'
import PropTypes from 'prop-types'

const Form = ({ onSubmit, children }) => (
    <form onSubmit={onSubmit} className={'mb-5'}>
        {children}
    </form>
)

Form.propTypes = {
    onSubmit: PropTypes.func,
    children: PropTypes.node,
}

export default Form
