export type ThemeTypes = {
    name: string;
    dark: boolean;
    variables?: object;
    colors: {
        primary?: string;
        secondary?: string;
        info?: string;
        success?: string;
        accent?: string;
        warning?: string;
        error?: string;
        light_primary?: string;
        light_secondary?: string;
        light_success?: string;
        light_error?: string;
        light_info?: string;
        light_warning?: string;
        textPrimary?: string;
        textSecondary?: string;
        borderColor?: string;
        hoverColor?: string;
        inputBorder?: string;
        containerBg?: string;
        surface?: string;
        'on-surface-variant'?: string;
        grey100?: string;
        grey200?: string;
        muted?:string;
        blue_light?:string
    };
};
