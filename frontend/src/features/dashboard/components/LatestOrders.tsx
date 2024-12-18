import React from 'react';

interface Order {
  id: number;
  from: string;
  to: string;
  weight: number;
  status: 'Paid' | 'Unpaid';
}

interface LatestOrdersProps {
  data: Order[];
}

const LatestOrders: React.FC<LatestOrdersProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="bg-white rounded-xl shadow-sm p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/2 mb-6"></div>
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center space-x-4">
                <div className="h-12 w-12 bg-gray-200 rounded"></div>
                <div className="flex-1">
                  <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                  <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl shadow-sm p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900">Latest Orders</h3>
        <button className="text-sm text-indigo-600 hover:text-indigo-900">
          View all
        </button>
      </div>

      <div className="space-y-6">
        {data.map((order) => (
          <div key={order.id} className="flex items-start space-x-4">
            <div className="flex-shrink-0">
              <div className="h-10 w-10 rounded-lg bg-gray-100 flex items-center justify-center">
                {/* Add package icon */}
              </div>
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium text-gray-900">
                  From {order.from}
                </p>
                <div className="ml-2">
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      order.status === 'Paid'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }`}
                  >
                    {order.status}
                  </span>
                </div>
              </div>
              <div className="mt-1">
                <p className="text-sm text-gray-500 truncate">
                  To {order.to} • {order.weight} kg
                </p>
              </div>
              <div className="mt-2 flex items-center">
                <button className="text-sm text-indigo-600 hover:text-indigo-900">
                  View details
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LatestOrders;
