import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/layout/Layout';
import Home from './components/Home';
import SearchPage from './components/SearchPage';
import Login from './components/Login';
import Register from './components/Register';
import UploadVideo from './components/UploadVideo';
import VideoPage from './components/VideoPage';

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/upload" element={<UploadVideo />} />
          <Route path="/watch/:id" element={<VideoPage />} />
          <Route path="/v1/videos/id/:id" element={<VideoPage />} />
        </Routes>
      </Layout>
    </Router>
  );
};

export default App;
