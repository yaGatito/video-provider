import styled from 'styled-components';

export const PageShell = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.xl};
  width: 100%;
`;

export const PageHeading = styled.h1`
  font-family: ${({ theme }) => theme.fonts.heading};
  color: ${({ theme }) => theme.colors.textPrimary};
  margin: 0;
`;

export const SectionBlock = styled.section`
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
`;

export const SectionTitle = styled.h2`
  font-family: ${({ theme }) => theme.fonts.heading};
  color: ${({ theme }) => theme.colors.textPrimary};
  margin: 0;
`;

export const ContentGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: ${({ theme }) => theme.spacing.lg};
`;

export const FormShell = styled.form`
  width: min(100%, 560px);
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.md};
  box-shadow: ${({ theme }) => theme.shadows.sm};
  padding: ${({ theme }) => theme.spacing.xl};
  display: grid;
  gap: ${({ theme }) => theme.spacing.lg};
`;

export const FormField = styled.div`
  display: grid;
  gap: ${({ theme }) => theme.spacing.sm};
`;

export const Label = styled.label`
  font-weight: 700;
  color: ${({ theme }) => theme.colors.textPrimary};
`;

export const TextInput = styled.input`
  width: 100%;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.85rem 0.95rem;
  background: ${({ theme }) => theme.colors.surface};
  transition: border-color ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base};

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

export const TextAreaInput = styled.textarea`
  width: 100%;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: 0.85rem 0.95rem;
  min-height: 140px;
  resize: vertical;
  background: ${({ theme }) => theme.colors.surface};
  transition: border-color ${({ theme }) => theme.transitions.base}, box-shadow ${({ theme }) => theme.transitions.base};

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.brand};
    box-shadow: 0 0 0 3px rgba(15, 76, 129, 0.15);
  }
`;

export const Message = styled.p<{ $tone?: 'success' | 'error' | 'info' }>`
  width: 100%;
  border-radius: ${({ theme }) => theme.radius.sm};
  padding: ${({ theme }) => theme.spacing.md};
  background: ${({ theme, $tone }) =>
    $tone === 'success'
      ? theme.colors.successBg
      : $tone === 'error'
      ? theme.colors.errorBg
      : theme.colors.surfaceAlt};
  color: ${({ theme, $tone }) =>
    $tone === 'success'
      ? theme.colors.successText
      : $tone === 'error'
      ? theme.colors.errorText
      : theme.colors.textSecondary};
`;

export const CenteredPage = styled(PageShell)`
  justify-items: center;
`;
