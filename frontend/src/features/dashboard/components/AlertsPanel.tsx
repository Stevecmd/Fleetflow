import React from 'react';

const AlertsPanel: React.FC<{ data: { id: number, type: string, message: string, priority: string, timestamp: string }[] }> = ({ data }) => {
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Alerts</h2>
      <ul>
        {data.map(alert => (
          <li key={alert.id} className="mb-2">
            <p className="font-semibold">Type: {alert.type}</p>
            <p>Message: {alert.message}</p>
            <p>Priority: {alert.priority}</p>
            <p>Timestamp: {alert.timestamp}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default AlertsPanel;