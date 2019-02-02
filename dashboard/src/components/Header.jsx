import React from 'react'
import PropTypes from 'prop-types'

const Header = ({ title, pretitle, children }) => (
    <div className={'header'}>
        <div className={'row align-items-end'}>
            <div className={'col'}>
                <h6 className={'header-pretitle'}>
                    {pretitle}
                </h6>
                <h1 className={'header-title'}>
                    {title}
                </h1>
            </div>
            {children && (
                <div className={'col-auto'}>
                    {children}
                </div>
            )}
        </div>
    </div>
)

Header.propTypes = {
    title: PropTypes.string.isRequired,
    pretitle: PropTypes.string.isRequired,
}

export default Header
