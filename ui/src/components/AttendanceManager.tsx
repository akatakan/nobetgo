import React, { useState, useEffect } from 'react';
import { ClipboardCheck, Clock, ChevronLeft, ChevronRight, Loader2, CheckCircle2, AlertCircle } from 'lucide-react';
import { scheduleApi, attendanceApi, departmentApi, employeeApi, shiftTypeApi } from '../services/api';
import type { Schedule, Attendance, Department, Employee, ShiftType } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const AttendanceManager: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [schedules, setSchedules] = useState<Schedule[]>([]);
    const [attendances, setAttendances] = useState<Attendance[]>([]);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [shiftTypes, setShiftTypes] = useState<ShiftType[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState<number | null>(null);
    const [showLogForm, setShowLogForm] = useState<number | null>(null);
    const [logForm, setLogForm] = useState({ start: '08:00', end: '16:00', notes: '' });

    // Manual Entry State
    const [showManualForm, setShowManualForm] = useState(false);
    const [manualForm, setManualForm] = useState({
        employeeId: 0,
        date: new Date().toISOString().split('T')[0],
        shiftTypeId: 0,
        start: '08:00',
        end: '17:00',
        notes: ''
    });

    useEffect(() => {
        Promise.all([
            departmentApi.list(),
            employeeApi.list(),
            shiftTypeApi.list()
        ]).then(([deptRes, empRes, stRes]) => {
            setDepartments(deptRes.data);
            setEmployees(empRes.data);
            setShiftTypes(stRes.data);
            if (deptRes.data.length > 0 && selectedDept === 0) setSelectedDept(deptRes.data[0].ID);
        }).catch(console.error);
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [schedRes, attRes] = await Promise.all([
                scheduleApi.getMonthly(month, year, selectedDept || undefined),
                attendanceApi.getReport(month, year, selectedDept || undefined),
            ]);
            setSchedules(schedRes.data || []);
            setAttendances(attRes.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (selectedDept) fetchData();
    }, [month, year, selectedDept]);

    const prevMonth = () => {
        if (month === 1) { setMonth(12); setYear(year - 1); } else setMonth(month - 1);
    };
    const nextMonth = () => {
        if (month === 12) { setMonth(1); setYear(year + 1); } else setMonth(month + 1);
    };

    const attendanceMap = new Map<number, Attendance>();
    attendances.forEach(a => attendanceMap.set(a.ScheduleID, a));

    const handleLogAttendance = async (scheduleId: number) => {
        setSaving(scheduleId);
        try {
            const dateObj = schedules.find(s => s.ID === scheduleId)?.Date;
            const dateStr = dateObj ? new Date(dateObj).toISOString().split('T')[0] : new Date().toISOString().split('T')[0];

            let endDateStr = dateStr;
            if (logForm.end < logForm.start) {
                const d = dateObj ? new Date(dateObj) : new Date();
                d.setDate(d.getDate() + 1);
                endDateStr = d.toISOString().split('T')[0];
            }

            const payload = {
                schedule_id: scheduleId,
                actual_start_time: `${dateStr}T${logForm.start}:00Z`,
                actual_end_time: `${endDateStr}T${logForm.end}:00Z`,
                notes: logForm.notes,
            };

            const existing = attendanceMap.get(scheduleId);
            if (existing) {
                await attendanceApi.update(existing.ID, payload);
            } else {
                await attendanceApi.log(payload);
            }

            setShowLogForm(null);
            setLogForm({ start: '08:00', end: '16:00', notes: '' });
            fetchData();
        } catch (err) {
            alert('Puantaj kayıt hatası: ' + err);
        } finally {
            setSaving(null);
        }
    };

    // Group schedules by date
    const sortedSchedules = [...schedules].sort((a, b) => new Date(a.Date).getTime() - new Date(b.Date).getTime());
    const schedulesByDate: Record<string, Schedule[]> = {};
    sortedSchedules.forEach(s => {
        const key = new Date(s.Date).toISOString().split('T')[0];
        if (!schedulesByDate[key]) schedulesByDate[key] = [];
        schedulesByDate[key].push(s);
    });

    // Stats
    const totalSchedules = schedules.length;
    const loggedCount = schedules.filter(s => attendanceMap.has(s.ID)).length;
    const overtimeCount = attendances.filter(a => a.IsOvertime).length;
    const totalOvertimeHours = attendances.reduce((sum, a) => sum + (a.OvertimeHours || 0), 0);

    const handleManualEntry = async () => {
        if (!manualForm.employeeId || !manualForm.shiftTypeId || !manualForm.date) {
            alert('Lütfen personel, tarih ve vardiya tipini seçin.');
            return;
        }
        setSaving(-1); // Special loading state
        try {
            // 1. Create Schedule
            const schedRes = await scheduleApi.create({
                EmployeeID: manualForm.employeeId,
                ShiftTypeID: manualForm.shiftTypeId,
                Date: new Date(manualForm.date),
            });

            if (!schedRes.data || !schedRes.data.ID) throw new Error('Schedule creation failed');

            // 2. Log Attendance
            let endDateStr = manualForm.date;
            if (manualForm.end < manualForm.start) {
                const d = new Date(manualForm.date);
                d.setDate(d.getDate() + 1);
                endDateStr = d.toISOString().split('T')[0];
            }

            await attendanceApi.log({
                schedule_id: schedRes.data.ID,
                actual_start_time: `${manualForm.date}T${manualForm.start}:00Z`,
                actual_end_time: `${endDateStr}T${manualForm.end}:00Z`,
                notes: manualForm.notes || 'Manuel Ek',
            });

            setShowManualForm(false);
            setManualForm({
                employeeId: 0,
                date: new Date().toISOString().split('T')[0],
                shiftTypeId: shiftTypes[0]?.ID || 0,
                start: '08:00',
                end: '17:00',
                notes: ''
            });
            fetchData();
        } catch (err) {
            alert('Kayıt hatası: ' + err);
        } finally {
            setSaving(null);
        }
    };

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Puantaj Kaydı</h2>
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
                        <Clock className="w-3.5 h-3.5" /> Manuel Ekle
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

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                    <div className="text-xs text-gray-400 mb-1">Toplam Nöbet</div>
                    <div className="text-2xl font-bold text-blue-400">{totalSchedules}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-emerald-500/5 to-emerald-500/10 border-emerald-500/10">
                    <div className="text-xs text-gray-400 mb-1">Kaydedilen</div>
                    <div className="text-2xl font-bold text-emerald-400">{loggedCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-amber-500/5 to-amber-500/10 border-amber-500/10">
                    <div className="text-xs text-gray-400 mb-1">Ek Mesai Sayısı</div>
                    <div className="text-2xl font-bold text-amber-400">{overtimeCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-red-500/5 to-red-500/10 border-red-500/10">
                    <div className="text-xs text-gray-400 mb-1">Toplam Ek Mesai</div>
                    <div className="text-2xl font-bold text-red-400">{totalOvertimeHours.toFixed(1)} saat</div>
                </div>
            </div>

            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />
                    Yükleniyor...
                </div>
            ) : Object.keys(schedulesByDate).length === 0 ? (
                <div className="glass-card p-12 text-center text-gray-500">
                    <ClipboardCheck className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Bu ay için nöbet bulunamadı.</p>
                    <p className="text-sm mt-1">Önce nöbet oluşturun, sonra puantaj kaydı yapabilirsiniz.</p>
                </div>
            ) : (
                <div className="space-y-4">
                    {Object.entries(schedulesByDate).map(([dateStr, daySchedules]) => {
                        const dateObj = new Date(dateStr);
                        const dayName = dateObj.toLocaleDateString('tr-TR', { weekday: 'long' });
                        const dayNum = dateObj.getDate();
                        const isPast = dateObj < new Date(now.getFullYear(), now.getMonth(), now.getDate());

                        return (
                            <div key={dateStr} className="glass-card overflow-hidden animate-slide-up">
                                <div className="px-5 py-3 bg-white/[0.02] border-b border-white/5 flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-lg bg-blue-500/10 flex items-center justify-center font-bold text-blue-400">
                                        {dayNum}
                                    </div>
                                    <div>
                                        <div className="font-semibold text-sm">{dayName}</div>
                                        <div className="text-xs text-gray-500">{dateStr}</div>
                                    </div>
                                    {isPast && (
                                        <span className="ml-auto text-xs text-gray-600 bg-white/5 px-2 py-1 rounded-md">Geçmiş</span>
                                    )}
                                </div>
                                <div className="divide-y divide-white/[0.03]">
                                    {daySchedules.map((sched) => {
                                        const att = attendanceMap.get(sched.ID);
                                        const isLogged = !!att;
                                        const isFormOpen = showLogForm === sched.ID;

                                        return (
                                            <div key={sched.ID} className="px-5 py-3">
                                                <div className="flex items-center gap-4">
                                                    <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${isLogged ? 'bg-emerald-500/15' : 'bg-gray-500/10'}`}>
                                                        {isLogged
                                                            ? <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                                                            : <AlertCircle className="w-4 h-4 text-gray-500" />
                                                        }
                                                    </div>
                                                    <div className="flex-1">
                                                        <div className="font-medium text-sm">
                                                            {sched.Employee?.FirstName} {sched.Employee?.LastName}
                                                        </div>
                                                        <div className="text-xs text-gray-500">
                                                            {sched.ShiftType?.Name}
                                                            {isLogged && att && (
                                                                <span className="ml-2 text-emerald-400">
                                                                    ✓ {new Date(att.ActualStartTime).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' })} - {new Date(att.ActualEndTime).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' })}
                                                                </span>
                                                            )}
                                                        </div>
                                                    </div>
                                                    {isLogged && att?.IsOvertime && (
                                                        <span className="text-xs px-2 py-1 rounded-md bg-amber-500/10 text-amber-400 font-medium">
                                                            +{att.OvertimeHours?.toFixed(1)}s mesai
                                                        </span>
                                                    )}
                                                    {!isLogged ? (
                                                        <button
                                                            onClick={() => {
                                                                setShowLogForm(isFormOpen ? null : sched.ID);
                                                                setLogForm({ start: '08:00', end: '16:00', notes: '' });
                                                            }}
                                                            className="btn-ghost text-xs py-1.5 px-3"
                                                        >
                                                            <Clock className="w-3 h-3" />
                                                            Kaydet
                                                        </button>
                                                    ) : (
                                                        <button
                                                            onClick={() => {
                                                                if (att) {
                                                                    const start = new Date(att.ActualStartTime).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' });
                                                                    const end = new Date(att.ActualEndTime).toLocaleTimeString('tr-TR', { hour: '2-digit', minute: '2-digit' });
                                                                    setLogForm({ start, end, notes: att.Notes || '' });
                                                                    setShowLogForm(isFormOpen ? null : sched.ID);
                                                                }
                                                            }}
                                                            className="btn-ghost text-xs py-1.5 px-3 opacity-60 hover:opacity-100"
                                                        >
                                                            Düzenle
                                                        </button>
                                                    )}
                                                </div>

                                                {/* Inline Log Form */}
                                                {isFormOpen && (
                                                    <div className="mt-3 pl-12 flex items-end gap-3 animate-slide-down">
                                                        <div className="space-y-1">
                                                            <label className="text-[10px] text-gray-500">Giriş</label>
                                                            <input
                                                                type="time"
                                                                className="glass-input py-1.5 px-2 text-sm w-28"
                                                                value={logForm.start}
                                                                onChange={(e) => setLogForm({ ...logForm, start: e.target.value })}
                                                            />
                                                        </div>
                                                        <div className="space-y-1">
                                                            <label className="text-[10px] text-gray-500">Çıkış</label>
                                                            <input
                                                                type="time"
                                                                className="glass-input py-1.5 px-2 text-sm w-28"
                                                                value={logForm.end}
                                                                onChange={(e) => setLogForm({ ...logForm, end: e.target.value })}
                                                            />
                                                        </div>
                                                        <div className="space-y-1 flex-1">
                                                            <label className="text-[10px] text-gray-500">Not</label>
                                                            <input
                                                                placeholder="İsteğe bağlı"
                                                                className="glass-input py-1.5 px-2 text-sm w-full"
                                                                value={logForm.notes}
                                                                onChange={(e) => setLogForm({ ...logForm, notes: e.target.value })}
                                                            />
                                                        </div>
                                                        <button
                                                            onClick={() => handleLogAttendance(sched.ID)}
                                                            disabled={saving === sched.ID}
                                                            className="btn-success text-xs py-1.5 px-4"
                                                        >
                                                            {saving === sched.ID ? <Loader2 className="w-3 h-3 animate-spin" /> : <CheckCircle2 className="w-3 h-3" />}
                                                            Onayla
                                                        </button>
                                                        <button onClick={() => setShowLogForm(null)} className="btn-ghost text-xs py-1.5 px-2">
                                                            İptal
                                                        </button>
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
                        <h3 className="text-lg font-semibold mb-4">Manuel Puantaj Ekle</h3>
                        <div className="space-y-4">
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Personel</label>
                                <select
                                    className="glass-input w-full"
                                    value={manualForm.employeeId}
                                    onChange={(e) => setManualForm({ ...manualForm, employeeId: Number(e.target.value) })}
                                >
                                    <option value={0}>Personel Seçin</option>
                                    {employees
                                        .filter(e => selectedDept === 0 || e.DepartmentID === selectedDept)
                                        .map((e) => (
                                            <option key={e.ID} value={e.ID}>{e.FirstName} {e.LastName}</option>
                                        ))}
                                </select>
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Tarih</label>
                                <input
                                    type="date"
                                    className="glass-input w-full"
                                    value={manualForm.date}
                                    onChange={(e) => setManualForm({ ...manualForm, date: e.target.value })}
                                />
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Vardiya Tipi (Referans)</label>
                                <select
                                    className="glass-input w-full"
                                    value={manualForm.shiftTypeId}
                                    onChange={(e) => setManualForm({ ...manualForm, shiftTypeId: Number(e.target.value) })}
                                >
                                    <option value={0}>Seçin</option>
                                    {shiftTypes.map((st) => (
                                        <option key={st.ID} value={st.ID}>{st.Name}</option>
                                    ))}
                                </select>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Giriş</label>
                                    <input
                                        type="time"
                                        className="glass-input w-full"
                                        value={manualForm.start}
                                        onChange={(e) => setManualForm({ ...manualForm, start: e.target.value })}
                                    />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Çıkış</label>
                                    <input
                                        type="time"
                                        className="glass-input w-full"
                                        value={manualForm.end}
                                        onChange={(e) => setManualForm({ ...manualForm, end: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Not</label>
                                <input
                                    className="glass-input w-full"
                                    value={manualForm.notes}
                                    onChange={(e) => setManualForm({ ...manualForm, notes: e.target.value })}
                                />
                            </div>
                        </div>
                        <div className="flex justify-end gap-3 mt-6">
                            <button
                                onClick={() => setShowManualForm(false)}
                                className="btn-ghost"
                            >
                                İptal
                            </button>
                            <button
                                onClick={handleManualEntry}
                                disabled={saving !== null}
                                className="btn-primary"
                            >
                                {saving === -1 ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
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
