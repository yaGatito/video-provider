import styled from 'styled-components';

export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'ghost';
}

const Button = styled.button<ButtonProps>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
  padding: 0.7rem 1.1rem;
  border-radius: ${({ theme }) => theme.radius.sm};
  border: 1px solid transparent;
  font-weight: 600;
  cursor: pointer;
  transition: background ${({ theme }) => theme.transitions.base},
    color ${({ theme }) => theme.transitions.base},
    border-color ${({ theme }) => theme.transitions.base},
    transform ${({ theme }) => theme.transitions.base};

  ${({ variant, theme }) => {
    switch (variant) {
      case 'secondary':
        return `
          background: ${theme.colors.surfaceAlt};
          color: ${theme.colors.textPrimary};
          border-color: ${theme.colors.border};

          &:hover {
            background: ${theme.colors.surface};
          }
        `;
      case 'ghost':
        return `
          background: transparent;
          color: ${theme.colors.brand};
          border-color: transparent;

          &:hover {
            background: rgba(15, 76, 129, 0.08);
          }
        `;
      default:
        return `
          background: ${theme.colors.brand};
          color: white;

          &:hover {
            background: ${theme.colors.brandHover};
          }
        `;
    }
  }}

  &:active {
    transform: translateY(1px);
  }
`;

Button.defaultProps = {
  variant: 'primary',
};

export default Button;
