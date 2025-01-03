import React from 'react';

const FleetOverview: React.FC<{ data: { vehicleId: string, status: string, location: string, driver: string, lastUpdate: string }[] }> = ({ data }) => {
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Fleet Overview</h2>
      <ul>
        {data.map(vehicle => (
          <li key={vehicle.vehicleId} className="mb-2">
            <p className="font-semibold">Vehicle ID: {vehicle.vehicleId}</p>
            <p>Status: {vehicle.status}</p>
            <p>Location: {vehicle.location}</p>
            <p>Driver: {vehicle.driver}</p>
            <p>Last Update: {vehicle.lastUpdate}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default FleetOverview;