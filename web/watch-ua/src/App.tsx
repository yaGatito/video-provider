import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './components/Home';
import SearchPage from './components/SearchPage';
import Register from './components/Register'; // Import the new registration page

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/search" element={<SearchPage />} />
        <Route path="/register" element={<Register />} />
        {/* Add this route */}
      </Routes>
    </Router>
  );
};

export default App;