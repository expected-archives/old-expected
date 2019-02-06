import React, { ReactNode } from "react";

interface IColumn {
    title: string;
    key: string;
    render?: (data: object) => ReactNode;
}

interface IProps {
    title?: string;
    onRowClick?: () => void;
    dataSource: object[];
    columns: IColumn[];
}

export default ({ title, onRowClick, dataSource, columns }: IProps) => (
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
                            {columns.map(({ key, render }, index) => (
                                <td key={index}>{render ? render(data[key]) : data[key]}</td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
);
