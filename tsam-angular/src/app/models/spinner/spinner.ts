export declare type Size = 'default' | 'small';
export interface ISpinner {
    bdColor?: string;
    size?: Size;
    color?: string;
    type?: string;
    fullScreen?: boolean;
    zIndex?: number;
    loaderTemplate?: string;
    loadingTextTemplate?: string;
    autoDisableButton?: boolean;
}