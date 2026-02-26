import 'styled-components';

declare module 'styled-components' {
  export interface DefaultTheme {
    colors: {
      background: string;
      surface: string;
      surfaceAlt: string;
      textPrimary: string;
      textSecondary: string;
      muted: string;
      brand: string;
      brandHover: string;
      accent: string;
      border: string;
      successBg: string;
      successText: string;
      errorBg: string;
      errorText: string;
    };
    fonts: {
      heading: string;
      body: string;
    };
    spacing: {
      xs: string;
      sm: string;
      md: string;
      lg: string;
      xl: string;
      xxl: string;
    };
    radius: {
      sm: string;
      md: string;
      lg: string;
    };
    shadows: {
      sm: string;
      md: string;
    };
    breakpoints: {
      mobile: string;
      tablet: string;
    };
    transitions: {
      base: string;
    };
  }
}
