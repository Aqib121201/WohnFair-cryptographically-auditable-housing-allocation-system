'use client';

import { useState, useEffect } from 'react';
import { 
  BarChart3, 
  Users, 
  Clock, 
  TrendingUp, 
  Shield, 
  AlertTriangle,
  Plus,
  Search,
  Filter,
  Download,
  RefreshCw
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell } from 'recharts';

// Mock data for demonstration
const mockMetrics = {
  totalRequests: 1247,
  pendingRequests: 89,
  processedToday: 23,
  averageWaitTime: 14.2,
  fairnessScore: 94.7,
  systemEfficiency: 87.3
};

const mockQueueData = [
  { name: 'Refugees', count: 45, weight: 1.5, avgWait: 18.2 },
  { name: 'Disabled', count: 32, weight: 1.3, avgWait: 15.8 },
  { name: 'Seniors', count: 28, weight: 1.2, avgWait: 12.4 },
  { name: 'Low-Income', count: 67, weight: 1.1, avgWait: 10.1 },
  { name: 'General', count: 156, weight: 1.0, avgWait: 8.7 }
];

const mockFairnessTrend = [
  { date: '2024-01-01', score: 92.1, efficiency: 85.2 },
  { date: '2024-01-02', score: 93.4, efficiency: 86.1 },
  { date: '2024-01-03', score: 94.2, efficiency: 87.3 },
  { date: '2024-01-04', score: 94.7, efficiency: 87.3 },
  { date: '2024-01-05', score: 95.1, efficiency: 88.0 },
  { date: '2024-01-06', score: 94.8, efficiency: 87.8 },
  { date: '2024-01-07', score: 94.7, efficiency: 87.3 }
];

const mockRecentRequests = [
  { id: 'REQ-001', user: 'Ahmed K.', group: 'Refugee', urgency: 'High', submitted: '2 hours ago', status: 'Pending' },
  { id: 'REQ-002', user: 'Maria S.', group: 'Disabled', urgency: 'Medium', submitted: '3 hours ago', status: 'Processing' },
  { id: 'REQ-003', user: 'Hans M.', group: 'Senior', urgency: 'Low', submitted: '4 hours ago', status: 'Pending' },
  { id: 'REQ-004', user: 'Fatima A.', group: 'Low-Income', urgency: 'High', submitted: '5 hours ago', status: 'Pending' },
  { id: 'REQ-005', user: 'Klaus B.', group: 'General', urgency: 'Medium', submitted: '6 hours ago', status: 'Pending' }
];

const COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6'];

export default function DashboardPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedGroup, setSelectedGroup] = useState('all');

  const refreshData = async () => {
    setIsLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));
    setIsLoading(false);
  };

  useEffect(() => {
    refreshData();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Shield className="h-8 w-8 text-primary-600" />
              <span className="ml-2 text-xl font-bold text-gray-900">WohnFair Dashboard</span>
            </div>
            <div className="flex items-center space-x-4">
              <button
                onClick={refreshData}
                disabled={isLoading}
                className="btn btn-outline btn-sm"
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
                Refresh
              </button>
              <button className="btn btn-primary btn-sm">
                <Plus className="h-4 w-4 mr-2" />
                New Request
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Metrics Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6 mb-8">
          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-blue-100 rounded-lg">
                <Users className="h-6 w-6 text-blue-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Total Requests</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.totalRequests}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-yellow-100 rounded-lg">
                <Clock className="h-6 w-6 text-yellow-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Pending</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.pendingRequests}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-green-100 rounded-lg">
                <TrendingUp className="h-6 w-6 text-green-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Processed Today</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.processedToday}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-purple-100 rounded-lg">
                <Clock className="h-6 w-6 text-purple-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Avg Wait (days)</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.averageWaitTime}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-indigo-100 rounded-lg">
                <Shield className="h-6 w-6 text-indigo-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Fairness Score</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.fairnessScore}%</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center">
              <div className="p-2 bg-emerald-100 rounded-lg">
                <BarChart3 className="h-6 w-6 text-emerald-600" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Efficiency</p>
                <p className="text-2xl font-bold text-gray-900">{mockMetrics.systemEfficiency}%</p>
              </div>
            </div>
          </div>
        </div>

        {/* Charts Row */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Fairness Trend */}
          <div className="card p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Fairness & Efficiency Trend</h3>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={mockFairnessTrend}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="score" stroke="#3B82F6" strokeWidth={2} name="Fairness Score" />
                <Line type="monotone" dataKey="efficiency" stroke="#10B981" strokeWidth={2} name="System Efficiency" />
              </LineChart>
            </ResponsiveContainer>
          </div>

          {/* Queue Distribution */}
          <div className="card p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Queue Distribution by Group</h3>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={mockQueueData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="count"
                >
                  {mockQueueData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Queue Status Table */}
        <div className="card p-6 mb-8">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-semibold text-gray-900">Queue Status by Group</h3>
            <div className="flex items-center space-x-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search groups..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
              </div>
              <select
                value={selectedGroup}
                onChange={(e) => setSelectedGroup(e.target.value)}
                className="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="all">All Groups</option>
                <option value="refugees">Refugees</option>
                <option value="disabled">Disabled</option>
                <option value="seniors">Seniors</option>
                <option value="low-income">Low-Income</option>
                <option value="general">General</option>
              </select>
              <button className="btn btn-outline btn-sm">
                <Download className="h-4 w-4 mr-2" />
                Export
              </button>
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Group
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Queue Count
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Weight
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Avg Wait (days)
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {mockQueueData.map((group) => (
                  <tr key={group.name} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 h-8 w-8">
                          <div className="h-8 w-8 rounded-full bg-primary-100 flex items-center justify-center">
                            <Users className="h-4 w-4 text-primary-600" />
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">{group.name}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{group.count}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{group.weight}x</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{group.avgWait}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        group.avgWait > 15 ? 'bg-red-100 text-red-800' :
                        group.avgWait > 10 ? 'bg-yellow-100 text-red-800' :
                        'bg-green-100 text-green-800'
                      }`}>
                        {group.avgWait > 15 ? 'Critical' : group.avgWait > 10 ? 'Warning' : 'Good'}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Recent Requests */}
        <div className="card p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-lg font-semibold text-gray-900">Recent Housing Requests</h3>
            <button className="btn btn-outline btn-sm">View All</button>
          </div>

          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Request ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    User
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Group
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Urgency
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Submitted
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {mockRecentRequests.map((request) => (
                  <tr key={request.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {request.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {request.user}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        request.group === 'Refugee' ? 'bg-blue-100 text-blue-800' :
                        request.group === 'Disabled' ? 'bg-purple-100 text-purple-800' :
                        request.group === 'Senior' ? 'bg-green-100 text-green-800' :
                        request.group === 'Low-Income' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {request.group}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        request.urgency === 'High' ? 'bg-red-100 text-red-800' :
                        request.urgency === 'Medium' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-green-100 text-green-800'
                      }`}>
                        {request.urgency}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {request.submitted}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        request.status === 'Processing' ? 'bg-blue-100 text-blue-800' :
                        'bg-yellow-100 text-yellow-800'
                      }`}>
                        {request.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
