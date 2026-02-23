import React from 'react';
import './Profile.css';

interface User {
  username: string;
  email: string;
}

const Profile: React.FC<{ user: User }> = ({ user }) => {
  return (
    <div className="profile-container">
      <h1>User Profile</h1>
      <div className="profile-card">
        <img src="/default-profile.png" alt="User Profile" className="profile-image" />
        <div className="profile-info">
          <h2>{user.username}</h2>
          <p>Email: {user.email}</p>
          <p>Member since: January 2023</p>
        </div>
      </div>
    </div>
  );
};

export default Profile;
