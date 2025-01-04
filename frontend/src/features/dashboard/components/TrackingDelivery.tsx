import React from 'react';

interface Delivery {
  from: {
    code: string;
    time: string;
  };
  to: {
    code: string;
    time: string;
  };
}

/**
 * TrackingDelivery component displays the details of a delivery.
 *
 * @returns {ReactElement} A ReactElement representing the TrackingDelivery component.
 */
const TrackingDelivery: React.FC = () => {
  const delivery: Delivery = {
    from: {
      code: 'NYC',
      time: '12:30'
    },
    to: {
      code: 'SFO',
      time: '18:30'
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-sm p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900">Tracking Delivery</h3>
        <div className="text-sm text-gray-500">15 Deliveries</div>
      </div>

      <div className="bg-black rounded-lg p-6 text-white">
        <div className="flex justify-between items-center mb-8">
          <div>
            <div className="text-2xl font-bold">{delivery.from.code}</div>
            <div className="text-sm text-gray-400">{delivery.from.time}</div>
          </div>
          <div className="flex-1 mx-4">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="h-0.5 w-full bg-gray-700"></div>
              </div>
              <div className="relative flex justify-center">
                <button className="flex items-center justify-center w-12 h-12 rounded-full border-2 border-gray-700 bg-black text-white">
                  Start
                </button>
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold">{delivery.to.code}</div>
            <div className="text-sm text-gray-400">{delivery.to.time}</div>
          </div>
        </div>

        <div className="flex justify-between items-center text-sm">
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-green-400 rounded-full"></div>
            <span>On time</span>
          </div>
          <div className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-blue-400 rounded-full"></div>
            <span>4h 30m left</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TrackingDelivery;
