import React, { ReactNode } from "react";

interface IColumn {
    title: string;
    key: string;
    render?: (data: any) => ReactNode;
}

interface IProps<T> {
    title?: string;
    onRowClick?: (data: T) => any;
    dataSource: T[];
    columns: IColumn[];
}

export default <T,>({ title, onRowClick, dataSource, columns }: IProps<T>) => {
    const onClick = (data: T) => () => {
        if (onRowClick)
            onRowClick(data);
    };

    return (
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
                        {dataSource.map((data: any, index) => (
                            <tr key={index} onClick={onClick(data)}>
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
}
