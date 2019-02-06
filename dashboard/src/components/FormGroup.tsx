import React, { ReactNode } from "react";

interface IProps {
    name?: string;
    description?: string;
    children: ReactNode;
}

export default ({ name, description, children }: IProps) => (
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
);
