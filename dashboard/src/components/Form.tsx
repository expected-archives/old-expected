import React, { ReactNode, FormEvent } from "react";

interface IProps {
    onSubmit?: (event: FormEvent<HTMLFormElement>) => void;
    children: ReactNode;
}

export default ({ onSubmit, children }: IProps) => (
    <form onSubmit={onSubmit} className={'mb-5'}>
        {children}
    </form>
);
