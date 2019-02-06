import React, { ReactNode } from "react";

interface IProps {
    title?: string;
    children: ReactNode;
}

export default ({ title, children }: IProps) => (
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
);
