import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import Home from './components/Home';
import SearchPage from './components/SearchPage';
import Login from './components/Login';
import Profile from './components/Profile';
import Register from './components/Register';
import UploadVideo from './components/UploadVideo';
import VideoPage from './components/VideoPage';

const RequireAuth: React.FC<{ children: React.ReactElement }> = ({ children }) => {
  const token = localStorage.getItem('authToken');
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return children;
};

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/profile" element={<RequireAuth><Profile /></RequireAuth>} />
          <Route path="/upload" element={<UploadVideo />} />
          <Route path="/watch/:id" element={<VideoPage />} />
          <Route path="/videos/id/:id" element={<VideoPage />} />
        </Routes>
      </Layout>
    </Router>
  );
};

export default App;
