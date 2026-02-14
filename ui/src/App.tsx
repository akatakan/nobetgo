import React, { useState } from 'react';
import { Users, Clock, Calendar, BarChart3, Settings, Wand2 } from 'lucide-react';
import EmployeeManager from './components/EmployeeManager';
import ShiftTypeManager from './components/ShiftTypeManager';
import ScheduleWizard from './components/ScheduleWizard';
import ScheduleViewer from './components/ScheduleViewer';

const DashboardOverview: React.FC = () => (
  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
    <div className="bg-gray-700 p-6 rounded-xl border border-gray-600">
      <div className="text-gray-400 text-sm mb-1">Toplam Personel</div>
      <div className="text-3xl font-bold">12</div>
    </div>
    <div className="bg-gray-700 p-6 rounded-xl border border-gray-600">
      <div className="text-gray-400 text-sm mb-1">Bu Ayki Nöbetler</div>
      <div className="text-3xl font-bold">48</div>
    </div>
    <div className="bg-gray-700 p-6 rounded-xl border border-gray-600">
      <div className="text-gray-400 text-sm mb-1">Maliyet Tasarrufu</div>
      <div className="text-3xl font-bold text-green-400">%15</div>
    </div>
    <div className="col-span-1 md:col-span-3 bg-blue-900/20 p-8 rounded-2xl border border-blue-500/30 flex items-center justify-between">
      <div>
        <h3 className="text-xl font-bold text-blue-400 mb-2">Hızlı Başlangıç</h3>
        <p className="text-gray-400">Yeni bir çizelge oluşturmak için otomatik planlayıcıyı kullanın.</p>
      </div>
      <Wand2 className="w-12 h-12 text-blue-400 opacity-50" />
    </div>
  </div>
);

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState('dashboard');


  const navItems = [
    { id: 'dashboard', icon: BarChart3, label: 'Dashboard' },
    { id: 'schedule', icon: Calendar, label: 'Nöbet Listesi' },
    { id: 'employees', icon: Users, label: 'Personel' },
    { id: 'shifts', icon: Clock, label: 'Nöbet Tipleri' },
    { id: 'scheduler', icon: Wand2, label: 'Otomatik Planla' },
    { id: 'settings', icon: Settings, label: 'Ayarlar' },
  ];

  return (
    <div className="flex h-screen w-screen bg-gray-900 text-white overflow-hidden">
      {/* Sidebar */}
      <aside className="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
        <div className="p-6 font-bold text-xl text-blue-400 flex items-center gap-2">
          <Calendar className="w-8 h-8" />
          NöbetGo
        </div>
        <nav className="flex-1 px-4 py-4">
          <ul className="space-y-2">
            {navItems.map((item) => (
              <li key={item.id}>
                <button
                  onClick={() => setActiveTab(item.id)}
                  className={`w-full flex items-center gap-3 px-4 py-2 rounded-lg transition-colors ${activeTab === item.id
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-400 hover:bg-gray-700 hover:text-white'
                    }`}
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                </button>
              </li>
            ))}
          </ul>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-8">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">
            {navItems.find(i => i.id === activeTab)?.label}
          </h1>
          <p className="text-gray-400">Yönetim panelinize hoş geldiniz.</p>
        </header>

        <section className="bg-gray-800 rounded-xl border border-gray-700 p-6 shadow-xl">
          {activeTab === 'dashboard' && <DashboardOverview />}
          {activeTab === 'schedule' && <ScheduleViewer />}
          {activeTab === 'employees' && <EmployeeManager />}
          {activeTab === 'shifts' && <ShiftTypeManager />}
          {activeTab === 'scheduler' && <ScheduleWizard />}
        </section>
      </main>
    </div>
  );
};

export default App;
