import React from 'react'
import PropTypes from 'prop-types'

const TableCard = ({ title, onRowClick, dataSource, columns }) => (
    <div className={'card'}>
        {title && (
            <div className={'card-header'}>
                <h4>{title}</h4>
            </div>
        )}
        <div className={'card-table table-responsive'}>
            <table className={'table table-hover'}>
                <thead>
                    <tr>
                        {columns.map(({ title: columnTitle }, index) => (
                            <th key={index}>{columnTitle}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {dataSource.map((data, index) => (
                        <tr key={index} onClick={onRowClick && onRowClick.bind(this, data)}>
                            {columns.map(({ key, render }) => (
                                <td key={key}>{render ? render(data[key]) : data[key]}</td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
)

TableCard.propTypes = {
    title: PropTypes.string,
    onRowClick: PropTypes.func,
    dataSource: PropTypes.arrayOf(PropTypes.object).isRequired,
    columns: PropTypes.arrayOf(PropTypes.shape({
        title: PropTypes.string.isRequired,
        key: PropTypes.string.isRequired,
        render: PropTypes.func,
    })).isRequired,
}

export default TableCard
