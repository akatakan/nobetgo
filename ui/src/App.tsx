import React, { Suspense, lazy, useEffect, useState } from 'react';
import {
  Activity,
  Award,
  BarChart3,
  Building2,
  Calendar,
  CalendarOff,
  ClipboardCheck,
  Clock,
  FileBarChart,
  Lock as LockIcon,
  LogOut,
  Settings,
  ShieldCheck,
  TrendingUp,
  Users,
  Wand2,
} from 'lucide-react';
import Login from './components/Login';
import ChangePasswordModal from './components/ChangePasswordModal';
import { NotificationBell } from './components/NotificationBell';
import { departmentApi, employeeApi, scheduleApi, shiftTypeApi } from './services/api';
import './App.css';

const EmployeeManager = lazy(() => import('./components/EmployeeManager'));
const ShiftTypeManager = lazy(() => import('./components/ShiftTypeManager'));
const ScheduleWizard = lazy(() => import('./components/ScheduleWizard'));
const ScheduleViewer = lazy(() => import('./components/ScheduleViewer'));
const DepartmentManager = lazy(() => import('./components/DepartmentManager'));
const AttendanceManager = lazy(() => import('./components/AttendanceManager'));
const OvertimeReport = lazy(() => import('./components/OvertimeReport'));
const LeaveManager = lazy(() => import('./components/LeaveManager'));
const ApprovalManager = lazy(() => import('./components/ApprovalManager'));
const ReportingDashboard = lazy(() => import('./components/ReportingDashboard'));
const TitleManager = lazy(() => import('./components/TitleManager'));

const TabFallback: React.FC = () => (
  <div className="glass-card p-12 flex items-center justify-center text-gray-400 animate-fade-in">
    <div className="flex items-center gap-3">
      <div className="w-5 h-5 rounded-full border-2 border-blue-400/20 border-t-blue-400 animate-spin" />
      <span>Ekran yukleniyor...</span>
    </div>
  </div>
);

const DashboardOverview: React.FC<{ onNavigate: (tab: string) => void }> = ({ onNavigate }) => {
  const [stats, setStats] = useState({ employees: 0, shifts: 0, schedules: 0, departments: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [empRes, shiftRes, schedRes, deptRes] = await Promise.allSettled([
          employeeApi.list(),
          shiftTypeApi.list(),
          scheduleApi.getMonthly(new Date().getMonth() + 1, new Date().getFullYear()),
          departmentApi.list(),
        ]);
        setStats({
          employees: empRes.status === 'fulfilled' ? empRes.value.data.total : 0,
          shifts: shiftRes.status === 'fulfilled' ? shiftRes.value.data.length : 0,
          schedules: schedRes.status === 'fulfilled' ? schedRes.value.data.length : 0,
          departments: deptRes.status === 'fulfilled' ? deptRes.value.data.length : 0,
        });
      } catch {
        // Keep the dashboard resilient if one of the requests fails.
      }
      setLoading(false);
    };
    fetchStats();
  }, []);

  const statCards = [
    { label: 'Bolumler', value: stats.departments, icon: Building2, gradient: 'from-emerald-500/10 to-teal-500/10', borderColor: 'border-emerald-500/20', iconColor: 'text-emerald-400', valueColor: 'text-emerald-400' },
    { label: 'Toplam Personel', value: stats.employees, icon: Users, gradient: 'from-blue-500/10 to-cyan-500/10', borderColor: 'border-blue-500/20', iconColor: 'text-blue-400', valueColor: 'text-blue-400' },
    { label: 'Bu Ayki Nobetler', value: stats.schedules, icon: Calendar, gradient: 'from-purple-500/10 to-pink-500/10', borderColor: 'border-purple-500/20', iconColor: 'text-purple-400', valueColor: 'text-purple-400' },
    { label: 'Calisma Tipleri', value: stats.shifts, icon: Clock, gradient: 'from-amber-500/10 to-orange-500/10', borderColor: 'border-amber-500/20', iconColor: 'text-amber-400', valueColor: 'text-amber-400' },
  ];

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        {statCards.map((card, idx) => (
          <div
            key={card.label}
            className={`glass-card bg-gradient-to-br ${card.gradient} p-6 ${card.borderColor} transition-all duration-300 hover:scale-[1.02] cursor-default animate-slide-up`}
            style={{ animationDelay: `${idx * 100}ms`, animationFillMode: 'both' }}
          >
            <div className="flex items-center justify-between">
              <div>
                <div className="text-sm text-gray-400 mb-1">{card.label}</div>
                <div className={`text-3xl font-bold ${card.valueColor}`}>
                  {loading ? <div className="skeleton w-12 h-8 rounded" /> : card.value}
                </div>
              </div>
              <div className="w-12 h-12 rounded-xl flex items-center justify-center bg-white/5">
                <card.icon className={`w-6 h-6 ${card.iconColor}`} />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
        <div
          className="glass-card p-6 bg-gradient-to-br from-blue-500/5 to-indigo-500/5 border-blue-500/10 hover:border-blue-500/25 transition-all duration-300 cursor-pointer group animate-slide-up"
          style={{ animationDelay: '300ms', animationFillMode: 'both' }}
          onClick={() => onNavigate('attendance')}
        >
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-blue-500/20 to-indigo-500/20 flex items-center justify-center group-hover:scale-110 transition-transform">
              <ClipboardCheck className="w-7 h-7 text-blue-400" />
            </div>
            <div>
              <h3 className="font-bold text-lg text-white group-hover:text-blue-300 transition-colors">Puantaj Takibi</h3>
              <p className="text-sm text-gray-400">Otomatik giris-cikis ve mesai kaydi</p>
            </div>
            <TrendingUp className="w-5 h-5 text-gray-600 ml-auto group-hover:text-blue-400 group-hover:translate-x-1 transition-all" />
          </div>
        </div>

        <div
          className="glass-card p-6 bg-gradient-to-br from-emerald-500/5 to-teal-500/5 border-emerald-500/10 hover:border-emerald-500/25 transition-all duration-300 cursor-pointer group animate-slide-up"
          style={{ animationDelay: '400ms', animationFillMode: 'both' }}
          onClick={() => onNavigate('leaves')}
        >
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-teal-500/20 flex items-center justify-center group-hover:scale-110 transition-transform">
              <CalendarOff className="w-7 h-7 text-emerald-400" />
            </div>
            <div>
              <h3 className="font-bold text-lg text-white group-hover:text-emerald-300 transition-colors">Izin Yonetimi</h3>
              <p className="text-sm text-gray-400">Izin talep, onay ve bakiye takibi</p>
            </div>
            <Activity className="w-5 h-5 text-gray-600 ml-auto group-hover:text-emerald-400 group-hover:translate-x-1 transition-all" />
          </div>
        </div>

        <div
          className="glass-card p-6 bg-gradient-to-br from-purple-500/5 to-pink-500/5 border-purple-500/10 hover:border-purple-500/25 transition-all duration-300 cursor-pointer group animate-slide-up"
          style={{ animationDelay: '500ms', animationFillMode: 'both' }}
          onClick={() => onNavigate('reports')}
        >
          <div className="flex items-center gap-4">
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center group-hover:scale-110 transition-transform">
              <FileBarChart className="w-7 h-7 text-purple-400" />
            </div>
            <div>
              <h3 className="font-bold text-lg text-white group-hover:text-purple-300 transition-colors">Raporlar</h3>
              <p className="text-sm text-gray-400">Calisma saati ve trend analizi</p>
            </div>
            <BarChart3 className="w-5 h-5 text-gray-600 ml-auto group-hover:text-purple-400 group-hover:translate-x-1 transition-all" />
          </div>
        </div>
      </div>
    </div>
  );
};

const SettingsPage: React.FC = () => (
  <div className="space-y-8 animate-fade-in">
    <div>
      <div className="flex items-center gap-2 mb-4">
        <Award className="w-5 h-5 text-violet-400" />
        <h3 className="text-lg font-semibold text-gray-200">Unvan Yonetimi</h3>
      </div>
      <p className="text-sm text-gray-500 mb-4">Personel eklerken secilebilecek unvanlari buradan yonetebilirsiniz.</p>
      <Suspense fallback={<TabFallback />}>
        <TitleManager />
      </Suspense>
    </div>
  </div>
);

const SUBTITLES: Record<string, string> = {
  dashboard: 'Genel bakis ve hizli islemler',
  departments: 'Kat ve bolum tanimlarini yonetin',
  schedule: 'Aylik nobet takvimini goruntuleyin',
  employees: 'Personel kaydi olusturun ve yonetin',
  shifts: 'Calisma tiplerini tanimlayin',
  scheduler: 'Akilli algoritmalarla nobet cizelgesi olusturun',
  attendance: 'Otomatik giris-cikis ve mesai takibi',
  overtime: 'Fazla mesai hesaplama ve kural yonetimi',
  leaves: 'Izin talep, onay ve bakiye takibi',
  approvals: 'Onay bekleyen kayitlar ve denetim izi',
  reports: 'Calisma saatleri, izin ve trend analizleri',
  settings: 'Uygulama ayarlarini yonetin',
};

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [role, setRole] = useState<string | null>(localStorage.getItem('role'));
  const [showChangePassword, setShowChangePassword] = useState(false);

  const handleLoginSuccess = (nextToken: string, nextRole: string) => {
    setToken(nextToken);
    setRole(nextRole);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    setToken(null);
    setRole(null);
  };

  if (!token) {
    return <Login onLoginSuccess={handleLoginSuccess} />;
  }

  const navItems = [
    { id: 'dashboard', icon: BarChart3, label: 'Dashboard' },
    { id: 'departments', icon: Building2, label: 'Bolumler' },
    { id: 'schedule', icon: Calendar, label: 'Nobet Takvimi' },
    { id: 'employees', icon: Users, label: 'Personel' },
    { id: 'shifts', icon: Clock, label: 'Calisma Tipleri' },
    { id: 'scheduler', icon: Wand2, label: 'Otomatik Planla' },
    { id: 'attendance', icon: ClipboardCheck, label: 'Puantaj' },
    { id: 'leaves', icon: CalendarOff, label: 'Izinler' },
    { id: 'overtime', icon: TrendingUp, label: 'Mesai' },
    { id: 'approvals', icon: ShieldCheck, label: 'Onaylar' },
    { id: 'reports', icon: FileBarChart, label: 'Raporlar' },
    { id: 'settings', icon: Settings, label: 'Ayarlar' },
  ].filter((item) => {
    if (role !== 'admin') {
      return !['scheduler', 'approvals', 'settings', 'departments', 'shifts'].includes(item.id);
    }
    return true;
  });

  return (
    <div className="flex h-screen w-screen bg-[var(--bg-primary)] text-white overflow-hidden">
      <ChangePasswordModal isOpen={showChangePassword} onClose={() => setShowChangePassword(false)} />

      <aside className="w-64 bg-[var(--bg-secondary)] border-r border-white/[0.06] flex flex-col">
        <div className="p-6 flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
            <Calendar className="w-5 h-5 text-white" />
          </div>
          <span className="font-bold text-xl bg-gradient-to-r from-blue-400 to-blue-300 bg-clip-text text-transparent">
            NobetGo
          </span>
        </div>

        <nav className="flex-1 px-3 py-2 overflow-y-auto">
          <ul className="space-y-1">
            {navItems.map((item) => (
              <li key={item.id}>
                <button
                  onClick={() => setActiveTab(item.id)}
                  className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 text-sm font-medium relative ${
                    activeTab === item.id
                      ? 'bg-blue-500/10 text-blue-400'
                      : 'text-gray-400 hover:bg-white/[0.03] hover:text-gray-200'
                  }`}
                >
                  {activeTab === item.id && (
                    <div className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-5 bg-blue-500 rounded-r-full" />
                  )}
                  <item.icon className="w-[18px] h-[18px]" />
                  {item.label}
                </button>
              </li>
            ))}
          </ul>
        </nav>

        <div className="p-4 border-t border-white/[0.04]">
          <button
            onClick={() => setShowChangePassword(true)}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-gray-400 hover:bg-blue-500/10 hover:text-blue-400 transition-all text-sm font-medium mb-1"
          >
            <LockIcon className="w-[18px] h-[18px]" />
            Sifre Degistir
          </button>
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-gray-400 hover:bg-red-500/10 hover:text-red-400 transition-all text-sm font-medium mb-2"
          >
            <LogOut className="w-[18px] h-[18px]" />
            Guvenli Cikis
          </button>
          <div className="text-xs text-gray-600 text-center">NobetGo v0.3.1</div>
        </div>
      </aside>

      <main className="flex-1 overflow-auto">
        <header className="sticky top-0 z-10 bg-[var(--bg-primary)]/80 backdrop-blur-xl border-b border-white/[0.04] px-8 py-5 flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-white">{navItems.find((item) => item.id === activeTab)?.label}</h1>
            <p className="text-sm text-gray-500 mt-0.5">{SUBTITLES[activeTab]}</p>
          </div>
          <div className="flex items-center gap-4">
            <NotificationBell />
            <div className="flex items-center gap-2 px-3 py-1.5 bg-white/5 rounded-full border border-white/10">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-xs font-medium text-gray-300 uppercase tracking-wider">{role}</span>
            </div>
          </div>
        </header>

        <div className="p-8">
          {activeTab === 'dashboard' ? (
            <DashboardOverview onNavigate={setActiveTab} />
          ) : (
            <Suspense fallback={<TabFallback />}>
              {activeTab === 'departments' && <DepartmentManager />}
              {activeTab === 'schedule' && <ScheduleViewer />}
              {activeTab === 'employees' && <EmployeeManager />}
              {activeTab === 'shifts' && <ShiftTypeManager />}
              {activeTab === 'scheduler' && <ScheduleWizard onNavigate={setActiveTab} />}
              {activeTab === 'attendance' && <AttendanceManager />}
              {activeTab === 'leaves' && <LeaveManager />}
              {activeTab === 'overtime' && <OvertimeReport />}
              {activeTab === 'approvals' && <ApprovalManager />}
              {activeTab === 'reports' && <ReportingDashboard />}
              {activeTab === 'settings' && <SettingsPage />}
            </Suspense>
          )}
        </div>
      </main>
    </div>
  );
};

export default App;
