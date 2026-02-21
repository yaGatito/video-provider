import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './components/Home';
import SearchPage from './components/SearchPage';
import Register from './components/Register';
import UploadVideo from './components/UploadVideo';

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/register" element={<Register />} />
          <Route path="/upload" element={<UploadVideo />} />
        </Routes>
      </Layout>
    </Router>
  );
};

export default App;