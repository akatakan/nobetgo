import React, { useState, useEffect } from 'react';
import { ClipboardCheck, Clock, ChevronLeft, ChevronRight, Loader2, CheckCircle2, LogIn, LogOut, Plus, AlertCircle } from 'lucide-react';
import { timeEntryApi, departmentApi, employeeApi } from '../services/api';
import type { TimeEntry, Department, Employee } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const STATUS_LABELS: Record<string, { label: string; color: string }> = {
    pending: { label: 'Bekliyor', color: 'text-amber-400 bg-amber-500/10' },
    approved: { label: 'Onaylandı', color: 'text-emerald-400 bg-emerald-500/10' },
    rejected: { label: 'Reddedildi', color: 'text-red-400 bg-red-500/10' },
};

const ENTRY_TYPE_LABELS: Record<string, { label: string; color: string }> = {
    normal: { label: 'Normal', color: 'text-blue-400' },
    weekend: { label: 'Hafta Sonu', color: 'text-purple-400' },
    holiday: { label: 'Tatil', color: 'text-red-400' },
    overtime: { label: 'Fazla Mesai', color: 'text-amber-400' },
};

const AttendanceManager: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [entries, setEntries] = useState<TimeEntry[]>([]);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    // Manual entry form
    const [showManualForm, setShowManualForm] = useState(false);
    const [manualForm, setManualForm] = useState({
        employeeId: 0,
        clockIn: '',
        clockOut: '',
        breakMinutes: 0,
        notes: '',
    });

    useEffect(() => {
        Promise.all([
            departmentApi.list(),
            employeeApi.list(),
        ]).then(([deptRes, empRes]) => {
            setDepartments(deptRes.data);
            setEmployees(empRes.data);
            if (deptRes.data.length > 0 && selectedDept === 0) setSelectedDept(deptRes.data[0].ID);
        }).catch(console.error);
    }, []);

    const fetchEntries = async () => {
        setLoading(true);
        try {
            const startDate = `${year}-${String(month).padStart(2, '0')}-01`;
            const endMonth = month === 12 ? 1 : month + 1;
            const endYear = month === 12 ? year + 1 : year;
            const endDate = `${endYear}-${String(endMonth).padStart(2, '0')}-01`;

            const params: Record<string, string | number> = { start: startDate, end: endDate };
            if (selectedDept) params.department_id = selectedDept;

            const res = await timeEntryApi.list(params);
            setEntries(res.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (selectedDept) fetchEntries();
    }, [month, year, selectedDept]);

    const prevMonth = () => {
        if (month === 1) { setMonth(12); setYear(year - 1); } else setMonth(month - 1);
    };
    const nextMonth = () => {
        if (month === 12) { setMonth(1); setYear(year + 1); } else setMonth(month + 1);
    };

    // Clock In
    const handleClockIn = async (employeeId: number) => {
        setActionLoading(`clockin-${employeeId}`);
        try {
            await timeEntryApi.clockIn({ employee_id: employeeId });
            fetchEntries();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Giriş kaydı hatası');
        } finally {
            setActionLoading(null);
        }
    };

    // Clock Out
    const handleClockOut = async (employeeId: number) => {
        setActionLoading(`clockout-${employeeId}`);
        try {
            await timeEntryApi.clockOut({ employee_id: employeeId });
            fetchEntries();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Çıkış kaydı hatası');
        } finally {
            setActionLoading(null);
        }
    };

    // Manual Entry
    const handleManualEntry = async () => {
        if (!manualForm.employeeId || !manualForm.clockIn) {
            alert('Personel ve giriş saati gerekli.');
            return;
        }
        setActionLoading('manual');
        try {
            await timeEntryApi.create({
                employee_id: manualForm.employeeId,
                clock_in: manualForm.clockIn,
                clock_out: manualForm.clockOut || undefined,
                break_minutes: manualForm.breakMinutes,
                notes: manualForm.notes,
            });
            setShowManualForm(false);
            setManualForm({ employeeId: 0, clockIn: '', clockOut: '', breakMinutes: 0, notes: '' });
            fetchEntries();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Kayıt hatası');
        } finally {
            setActionLoading(null);
        }
    };

    // Delete
    const handleDelete = async (id: number) => {
        if (!confirm('Bu kaydı silmek istediğinize emin misiniz?')) return;
        try {
            await timeEntryApi.delete(id);
            fetchEntries();
        } catch (err) {
            console.error(err);
        }
    };

    // Group entries by date
    const entriesByDate: Record<string, TimeEntry[]> = {};
    [...entries]
        .sort((a, b) => new Date(a.clock_in).getTime() - new Date(b.clock_in).getTime())
        .forEach(e => {
            const key = new Date(e.clock_in).toISOString().split('T')[0];
            if (!entriesByDate[key]) entriesByDate[key] = [];
            entriesByDate[key].push(e);
        });

    // Stats
    const totalEntries = entries.length;
    const approvedCount = entries.filter(e => e.status === 'approved').length;
    const pendingCount = entries.filter(e => e.status === 'pending').length;
    const activeClockIns = entries.filter(e => !e.clock_out).length;

    // Find employees with active clock-in (no clock-out)
    const deptEmployees = employees.filter(e => !selectedDept || e.DepartmentID === selectedDept);
    const activeEmployeeIds = new Set(entries.filter(e => !e.clock_out).map(e => e.employee_id));

    const formatTime = (dateStr: string) =>
        new Date(dateStr).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' });

    const formatHours = (entry: TimeEntry) => {
        if (!entry.clock_out) return '—';
        const ms = new Date(entry.clock_out).getTime() - new Date(entry.clock_in).getTime();
        const hours = ms / 3600000 - (entry.break_minutes || 0) / 60;
        return `${Math.max(0, hours).toFixed(1)}s`;
    };

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Puantaj Takibi</h2>
                    <select
                        className="glass-input py-1.5 px-3 text-sm"
                        value={selectedDept}
                        onChange={(e) => setSelectedDept(Number(e.target.value))}
                    >
                        <option value={0} disabled>Bölüm Seçin</option>
                        {departments.map((d) => (
                            <option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>
                        ))}
                    </select>
                </div>
                <div className="flex items-center gap-2">
                    <button
                        onClick={() => setShowManualForm(true)}
                        className="btn-primary text-xs py-2 px-3 flex items-center gap-2"
                    >
                        <Plus className="w-3.5 h-3.5" /> Manuel Kayıt
                    </button>
                    <div className="h-6 w-px bg-white/10 mx-1"></div>
                    <button onClick={prevMonth} className="btn-ghost p-2">
                        <ChevronLeft className="w-4 h-4" />
                    </button>
                    <span className="min-w-[160px] text-center font-semibold text-lg">
                        {MONTHS_TR[month - 1]} {year}
                    </span>
                    <button onClick={nextMonth} className="btn-ghost p-2">
                        <ChevronRight className="w-4 h-4" />
                    </button>
                </div>
            </div>

            {/* Quick Clock In/Out */}
            <div className="glass-card p-5 bg-gradient-to-br from-blue-500/5 to-indigo-500/5 border-blue-500/10">
                <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2">
                    <Clock className="w-4 h-4 text-blue-400" /> Hızlı Giriş / Çıkış
                </h3>
                <div className="flex flex-wrap gap-2">
                    {deptEmployees.map(emp => {
                        const isActive = activeEmployeeIds.has(emp.ID);
                        const isLoading = actionLoading === `clockin-${emp.ID}` || actionLoading === `clockout-${emp.ID}`;
                        return (
                            <button
                                key={emp.ID}
                                onClick={() => isActive ? handleClockOut(emp.ID) : handleClockIn(emp.ID)}
                                disabled={isLoading}
                                className={`flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-medium transition-all duration-200 border ${isActive
                                    ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400 hover:bg-red-500/10 hover:border-red-500/20 hover:text-red-400'
                                    : 'bg-white/[0.03] border-white/10 text-gray-400 hover:bg-blue-500/10 hover:border-blue-500/20 hover:text-blue-400'
                                    }`}
                            >
                                {isLoading ? (
                                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                                ) : isActive ? (
                                    <LogOut className="w-3.5 h-3.5" />
                                ) : (
                                    <LogIn className="w-3.5 h-3.5" />
                                )}
                                {emp.FirstName} {emp.LastName}
                                {isActive && <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />}
                            </button>
                        );
                    })}
                    {deptEmployees.length === 0 && (
                        <p className="text-xs text-gray-600">Bu bölümde personel yok.</p>
                    )}
                </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                    <div className="text-xs text-gray-400 mb-1">Toplam Kayıt</div>
                    <div className="text-2xl font-bold text-blue-400">{totalEntries}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-emerald-500/5 to-emerald-500/10 border-emerald-500/10">
                    <div className="text-xs text-gray-400 mb-1">Onaylanan</div>
                    <div className="text-2xl font-bold text-emerald-400">{approvedCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-amber-500/5 to-amber-500/10 border-amber-500/10">
                    <div className="text-xs text-gray-400 mb-1">Bekleyen</div>
                    <div className="text-2xl font-bold text-amber-400">{pendingCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-purple-500/5 to-purple-500/10 border-purple-500/10">
                    <div className="text-xs text-gray-400 mb-1">Aktif Giriş</div>
                    <div className="text-2xl font-bold text-purple-400">{activeClockIns}</div>
                </div>
            </div>

            {/* Entries List */}
            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />
                    Yükleniyor...
                </div>
            ) : Object.keys(entriesByDate).length === 0 ? (
                <div className="glass-card p-12 text-center text-gray-500">
                    <ClipboardCheck className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Bu ay için puantaj kaydı bulunamadı.</p>
                    <p className="text-sm mt-1">Giriş/Çıkış butonlarını veya Manuel Kayıt'ı kullanarak kayıt oluşturun.</p>
                </div>
            ) : (
                <div className="space-y-4">
                    {Object.entries(entriesByDate).map(([dateStr, dayEntries]) => {
                        const dateObj = new Date(dateStr);
                        const dayName = dateObj.toLocaleDateString('tr-TR', { weekday: 'long' });
                        const dayNum = dateObj.getDate();
                        const isWeekend = dateObj.getDay() === 0 || dateObj.getDay() === 6;

                        return (
                            <div key={dateStr} className="glass-card overflow-hidden animate-slide-up">
                                <div className={`px-5 py-3 border-b border-white/5 flex items-center gap-3 ${isWeekend ? 'bg-purple-500/5' : 'bg-white/[0.02]'}`}>
                                    <div className={`w-10 h-10 rounded-lg flex items-center justify-center font-bold ${isWeekend ? 'bg-purple-500/10 text-purple-400' : 'bg-blue-500/10 text-blue-400'}`}>
                                        {dayNum}
                                    </div>
                                    <div>
                                        <div className="font-semibold text-sm">{dayName}</div>
                                        <div className="text-xs text-gray-500">{dateStr}</div>
                                    </div>
                                    {isWeekend && (
                                        <span className="ml-auto text-xs text-purple-400 bg-purple-500/10 px-2 py-1 rounded-md">Hafta Sonu</span>
                                    )}
                                </div>
                                <div className="divide-y divide-white/[0.03]">
                                    {dayEntries.map((entry) => {
                                        const statusInfo = STATUS_LABELS[entry.status] || STATUS_LABELS.pending;
                                        const typeInfo = ENTRY_TYPE_LABELS[entry.entry_type] || ENTRY_TYPE_LABELS.normal;

                                        return (
                                            <div key={entry.ID} className="px-5 py-3">
                                                <div className="flex items-center gap-4">
                                                    <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${entry.clock_out ? 'bg-emerald-500/15' : 'bg-amber-500/15'}`}>
                                                        {entry.clock_out
                                                            ? <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                                                            : <Clock className="w-4 h-4 text-amber-400 animate-pulse" />
                                                        }
                                                    </div>
                                                    <div className="flex-1">
                                                        <div className="font-medium text-sm">
                                                            {entry.employee?.FirstName} {entry.employee?.LastName}
                                                        </div>
                                                        <div className="text-xs text-gray-500 flex items-center gap-2">
                                                            <span className="text-emerald-400">{formatTime(entry.clock_in)}</span>
                                                            {entry.clock_out ? (
                                                                <>
                                                                    <span className="text-gray-600">→</span>
                                                                    <span className="text-blue-400">{formatTime(entry.clock_out)}</span>
                                                                    <span className="text-gray-500">({formatHours(entry)})</span>
                                                                </>
                                                            ) : (
                                                                <span className="text-amber-400 animate-pulse">devam ediyor...</span>
                                                            )}
                                                            {entry.break_minutes > 0 && (
                                                                <span className="text-gray-600">· {entry.break_minutes}dk mola</span>
                                                            )}
                                                        </div>
                                                    </div>
                                                    <span className={`text-[10px] px-2 py-0.5 rounded-md font-medium ${typeInfo.color}`}>
                                                        {typeInfo.label}
                                                    </span>
                                                    <span className={`text-[10px] px-2 py-0.5 rounded-md font-medium ${statusInfo.color}`}>
                                                        {statusInfo.label}
                                                    </span>
                                                    <span className="text-[10px] text-gray-600 bg-white/5 px-2 py-0.5 rounded-md">
                                                        {entry.source}
                                                    </span>
                                                    <button
                                                        onClick={() => handleDelete(entry.ID)}
                                                        className="btn-ghost text-xs py-1 px-2 text-red-400/60 hover:text-red-400"
                                                    >
                                                        Sil
                                                    </button>
                                                </div>
                                                {entry.notes && (
                                                    <div className="mt-1 pl-12 text-xs text-gray-500 italic">
                                                        {entry.notes}
                                                    </div>
                                                )}
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}

            {/* Manual Entry Modal */}
            {showManualForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 animate-fade-in">
                    <div className="glass-card p-6 w-full max-w-md bg-[#1e293b]">
                        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <AlertCircle className="w-5 h-5 text-blue-400" /> Manuel Puantaj Kaydı
                        </h3>
                        <div className="space-y-4">
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Personel</label>
                                <select
                                    className="glass-input w-full"
                                    value={manualForm.employeeId}
                                    onChange={(e) => setManualForm({ ...manualForm, employeeId: Number(e.target.value) })}
                                >
                                    <option value={0}>Personel Seçin</option>
                                    {deptEmployees.map((e) => (
                                        <option key={e.ID} value={e.ID}>{e.FirstName} {e.LastName}</option>
                                    ))}
                                </select>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Giriş Zamanı</label>
                                    <input
                                        type="datetime-local"
                                        className="glass-input w-full"
                                        value={manualForm.clockIn}
                                        onChange={(e) => setManualForm({ ...manualForm, clockIn: e.target.value })}
                                    />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Çıkış Zamanı</label>
                                    <input
                                        type="datetime-local"
                                        className="glass-input w-full"
                                        value={manualForm.clockOut}
                                        onChange={(e) => setManualForm({ ...manualForm, clockOut: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Mola (dakika)</label>
                                    <input
                                        type="number"
                                        className="glass-input w-full"
                                        value={manualForm.breakMinutes}
                                        onChange={(e) => setManualForm({ ...manualForm, breakMinutes: Number(e.target.value) })}
                                    />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Not</label>
                                    <input
                                        className="glass-input w-full"
                                        placeholder="İsteğe bağlı"
                                        value={manualForm.notes}
                                        onChange={(e) => setManualForm({ ...manualForm, notes: e.target.value })}
                                    />
                                </div>
                            </div>
                        </div>
                        <div className="flex justify-end gap-3 mt-6">
                            <button onClick={() => setShowManualForm(false)} className="btn-ghost">
                                İptal
                            </button>
                            <button
                                onClick={handleManualEntry}
                                disabled={actionLoading === 'manual'}
                                className="btn-primary"
                            >
                                {actionLoading === 'manual' ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                                Kaydet
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AttendanceManager;
