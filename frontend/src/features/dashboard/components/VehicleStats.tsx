import React from 'react';
import { VehicleData } from '../dashboardSlice'; // Correcting the import path

interface VehicleStatsProps {
  data: VehicleData | null;
}

const VehicleStats: React.FC<VehicleStatsProps> = ({ data }) => {
  if (!data) {
    return (
      <div className="bg-white rounded-xl shadow-sm p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-3/4 mb-6"></div>
          <div className="space-y-6">
            <div>
              <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
              <div className="h-8 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm p-6">
      <h2 className="text-2xl font-bold mb-4">Vehicle Stats</h2>
      <p>Plate Number: {data.plate_number}</p>
      <p>Type: {data.type}</p>
      <p>Make: {data.make}</p>
      <p>Model: {data.model}</p>
      <p>Year: {data.year}</p>
      <p>Capacity: {data.capacity}</p>
      <p>Fuel Type: {data.fuel_type}</p>
      <p>Status ID: {data.status_id}</p>
      <p>Mileage: {data.mileage}</p>
      <p>Created At: {data.created_at}</p>
      <p>Updated At: {data.updated_at}</p>
    </div>
  );
};

export default VehicleStats;
