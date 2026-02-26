import React from 'react';
import { render, screen } from '@testing-library/react';
import axios from 'axios';
import { ThemeProvider } from 'styled-components';
import App from './App';
import { theme } from './styles/theme';
import { GlobalStyles } from './styles/GlobalStyles';

jest.mock('axios');

const mockedAxios = axios as jest.Mocked<typeof axios>;

const renderAt = (path: string) => {
  window.history.pushState({}, 'test', path);
  return render(
    <ThemeProvider theme={theme}>
      <GlobalStyles />
      <App />
    </ThemeProvider>
  );
};

describe('App smoke tests', () => {
  beforeEach(() => {
    mockedAxios.get.mockResolvedValue({ data: [] } as never);
    mockedAxios.post.mockResolvedValue({ data: {} } as never);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('renders layout with header and footer', () => {
    renderAt('/search');

    expect(screen.getAllByText('Watch UA').length).toBeGreaterThan(0);
    expect(screen.getByText(/All rights reserved\./i)).toBeInTheDocument();
  });

  test('renders Home route', async () => {
    renderAt('/');

    expect(await screen.findByText('Featured Video')).toBeInTheDocument();
  });

  test('renders Search route', () => {
    renderAt('/search');

    expect(screen.getByRole('heading', { name: 'Search' })).toBeInTheDocument();
  });

  test('renders Login route', () => {
    renderAt('/login');

    expect(screen.getByRole('heading', { name: 'Login' })).toBeInTheDocument();
  });

  test('renders Register route', () => {
    renderAt('/register');

    expect(screen.getByRole('heading', { name: 'Register' })).toBeInTheDocument();
  });

  test('renders Upload route', () => {
    renderAt('/upload');

    expect(screen.getByRole('heading', { name: 'Upload Video' })).toBeInTheDocument();
  });
});
