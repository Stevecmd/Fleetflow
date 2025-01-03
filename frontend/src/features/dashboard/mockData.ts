export const dashboardData = {
  stats: {
    activeDeliveries: 28,
    availableDrivers: 15,
    pendingOrders: 42,
    totalDeliveries: 150,
    alerts: 3,
    fleetUtilization: 85,
    warehouseCapacity: 72
  },

  deliveryTrends: {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    data: [23, 34, 28, 42, 38, 25, 30]
  },

  fleetOverview: [
    {
      vehicleId: "TRK001",
      status: "On Route",
      location: "Seattle, WA",
      driver: "John Smith",
      lastUpdate: "10 mins ago"
    },
    {
      vehicleId: "TRK002",
      status: "Available",
      location: "Portland, OR",
      driver: "Sarah Johnson",
      lastUpdate: "5 mins ago"
    }
  ],

  driverStatus: [
    {
      id: 1,
      name: "Alex Jackson",
      status: "Active",
      currentDelivery: "DEL001",
      hoursLogged: 6.5
    },
    {
      id: 2,
      name: "Maria Garcia",
      status: "Break",
      currentDelivery: null,
      hoursLogged: 4.0
    }
  ],

  alerts: [
    {
      id: 1,
      type: "Maintenance",
      message: "Vehicle TRK003 due for service",
      priority: "high",
      timestamp: "2024-03-01T10:00:00Z"
    },
    {
      id: 2,
      type: "Delay",
      message: "Delivery DEL005 running 30 mins late",
      priority: "medium",
      timestamp: "2024-03-01T09:30:00Z"
    }
  ]
};
