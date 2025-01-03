import React from 'react';

const DriverStatusList: React.FC<{ data: { id: number, name: string, status: string, currentDelivery: string | null, hoursLogged: number }[] }> = ({ data }) => {
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Driver Status</h2>
      <ul>
        {data.map(driver => (
          <li key={driver.id} className="mb-2">
            <p className="font-semibold">{driver.name}</p>
            <p>Status: {driver.status}</p>
            <p>Current Delivery: {driver.currentDelivery || 'None'}</p>
            <p>Hours Logged: {driver.hoursLogged}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default DriverStatusList;