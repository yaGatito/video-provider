import { DefaultTheme } from 'styled-components';

export const theme: DefaultTheme = {
  colors: {
    background: '#f4f7fb',
    surface: '#ffffff',
    surfaceAlt: '#eef3fa',
    textPrimary: '#1f2937',
    textSecondary: '#4b5563',
    muted: '#9ca3af',
    brand: '#0f4c81',
    brandHover: '#0b3a63',
    accent: '#e85d04',
    border: '#dbe3ef',
    successBg: '#def7e8',
    successText: '#166534',
    errorBg: '#fde7e7',
    errorText: '#991b1b'
  },
  fonts: {
    heading: "'Avenir Next', 'Segoe UI', sans-serif",
    body: "'Source Sans 3', 'Segoe UI', sans-serif"
  },
  spacing: {
    xs: '0.25rem',
    sm: '0.5rem',
    md: '1rem',
    lg: '1.5rem',
    xl: '2rem',
    xxl: '3rem'
  },
  radius: {
    sm: '6px',
    md: '10px',
    lg: '16px'
  },
  shadows: {
    sm: '0 4px 14px rgba(12, 43, 68, 0.08)',
    md: '0 12px 28px rgba(12, 43, 68, 0.12)'
  },
  breakpoints: {
    mobile: '640px',
    tablet: '960px'
  },
  transitions: {
    base: '0.2s ease'
  }
};
