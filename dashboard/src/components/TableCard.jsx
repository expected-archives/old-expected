import React from 'react'
import PropTypes from 'prop-types'

const TableCard = ({ title, dataSource, columns }) => (
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
                        {columns.map(({ title: columnTitle }) => (
                            <th key={columnTitle}>{columnTitle}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {dataSource.map(data => (
                        <tr>
                            {columns.map(({ key, render }) => (
                                <td>{render ? render(data[key]) : data[key]}</td>
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
    dataSource: PropTypes.arrayOf(PropTypes.object).isRequired,
    columns: PropTypes.arrayOf(PropTypes.shape({
        title: PropTypes.string.isRequired,
        key: PropTypes.string.isRequired,
        render: PropTypes.func,
    })).isRequired,
}

export default TableCard
