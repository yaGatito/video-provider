import React from 'react';
import './Profile.css';

const Profile: React.FC = () => {
  return (
    <div className="profile-container">
      <h1>User Profile</h1>
      <div className="profile-card">
        <img src="/default-profile.png" alt="User Profile" className="profile-image" />
        <div className="profile-info">
          <h2>John Doe</h2>
          <p>Email: john.doe@example.com</p>
          <p>Member since: January 2023</p>
        </div>
      </div>
    </div>
  );
};

export default Profile;