import React, { useState, useEffect } from 'react';
import { TrendingUp, ChevronLeft, ChevronRight, Loader2, DollarSign, Clock, Users } from 'lucide-react';
import { scheduleApi, attendanceApi, departmentApi, employeeApi } from '../services/api';
import type { Schedule, Attendance, Department, Employee } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

interface EmployeeStats {
    employee: Employee;
    totalShifts: number;
    totalHours: number;
    overtimeHours: number;
    normalCost: number;
    overtimeCost: number;
    totalCost: number;
}

const OvertimeReport: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [departments, setDepartments] = useState<Department[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [schedules, setSchedules] = useState<Schedule[]>([]);
    const [attendances, setAttendances] = useState<Attendance[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState(true);
    const [overtimeMultiplier] = useState(1.5);

    useEffect(() => {
        departmentApi.list().then(res => {
            setDepartments(res.data);
            if (res.data.length > 0 && selectedDept === 0) setSelectedDept(res.data[0].ID);
        }).catch(console.error);
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [schedRes, attRes, empRes] = await Promise.all([
                scheduleApi.getMonthly(month, year, selectedDept || undefined),
                attendanceApi.getReport(month, year, selectedDept || undefined),
                employeeApi.list(),
            ]);
            setSchedules(schedRes.data || []);
            setAttendances(attRes.data || []);
            setEmployees(empRes.data || []);
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

    // Build per-employee stats
    const attendanceMap = new Map<number, Attendance>();
    attendances.forEach(a => attendanceMap.set(a.ScheduleID, a));

    const deptEmployees = employees.filter(e => !selectedDept || e.DepartmentID === selectedDept);

    const empStatsMap = new Map<number, EmployeeStats>();
    deptEmployees.forEach(emp => {
        empStatsMap.set(emp.ID, {
            employee: emp,
            totalShifts: 0,
            totalHours: 0,
            overtimeHours: 0,
            normalCost: 0,
            overtimeCost: 0,
            totalCost: 0,
        });
    });

    schedules.forEach(sched => {
        const stat = empStatsMap.get(sched.EmployeeID);
        if (!stat) return;
        stat.totalShifts++;

        const att = attendanceMap.get(sched.ID);
        if (att) {
            const hours = (new Date(att.ActualEndTime).getTime() - new Date(att.ActualStartTime).getTime()) / 3600000;
            const normalHours = Math.min(hours, 8);
            const overtime = Math.max(0, hours - 8);

            stat.totalHours += hours;
            stat.overtimeHours += overtime;
            stat.normalCost += normalHours * stat.employee.HourlyRate;
            stat.overtimeCost += overtime * stat.employee.HourlyRate * overtimeMultiplier;
        } else {
            // No attendance logged — assume standard 8h
            stat.totalHours += 8;
            stat.normalCost += 8 * stat.employee.HourlyRate;
        }
        stat.totalCost = stat.normalCost + stat.overtimeCost;
    });

    const empStats = Array.from(empStatsMap.values()).filter(s => s.totalShifts > 0).sort((a, b) => b.totalCost - a.totalCost);

    // Totals
    const totalNormalCost = empStats.reduce((s, e) => s + e.normalCost, 0);
    const totalOvertimeCost = empStats.reduce((s, e) => s + e.overtimeCost, 0);
    const totalCost = totalNormalCost + totalOvertimeCost;
    const totalOvertimeHours = empStats.reduce((s, e) => s + e.overtimeHours, 0);

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Ek Mesai & Maliyet Raporu</h2>
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

            {/* Summary Cards */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><DollarSign className="w-3 h-3" /> Normal Maliyet</div>
                    <div className="text-2xl font-bold text-blue-400">₺{totalNormalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-amber-500/5 to-amber-500/10 border-amber-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><TrendingUp className="w-3 h-3" /> Ek Mesai Maliyeti</div>
                    <div className="text-2xl font-bold text-amber-400">₺{totalOvertimeCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-emerald-500/5 to-emerald-500/10 border-emerald-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><DollarSign className="w-3 h-3" /> Toplam Maliyet</div>
                    <div className="text-2xl font-bold text-emerald-400">₺{totalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-red-500/5 to-red-500/10 border-red-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><Clock className="w-3 h-3" /> Ek Mesai Saati</div>
                    <div className="text-2xl font-bold text-red-400">{totalOvertimeHours.toFixed(1)} saat</div>
                </div>
            </div>

            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />
                    Yükleniyor...
                </div>
            ) : empStats.length === 0 ? (
                <div className="glass-card p-12 text-center text-gray-500">
                    <Users className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Bu dönem için veri bulunamadı.</p>
                </div>
            ) : (
                <div className="glass-card overflow-hidden">
                    {/* Table Header */}
                    <div className="grid grid-cols-8 gap-2 px-5 py-3 bg-white/[0.02] border-b border-white/5 text-xs font-semibold text-gray-400 uppercase">
                        <div className="col-span-2">Personel</div>
                        <div className="text-center">Nöbet</div>
                        <div className="text-center">Toplam Saat</div>
                        <div className="text-center">Ek Mesai</div>
                        <div className="text-right">Normal ₺</div>
                        <div className="text-right">Mesai ₺</div>
                        <div className="text-right">Toplam ₺</div>
                    </div>

                    {/* Rows */}
                    {empStats.map((stat, idx) => (
                        <div
                            key={stat.employee.ID}
                            className="grid grid-cols-8 gap-2 px-5 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors animate-slide-up items-center"
                            style={{ animationDelay: `${idx * 30}ms`, animationFillMode: 'both' }}
                        >
                            <div className="col-span-2 flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/15 to-purple-500/15 flex items-center justify-center text-xs font-bold text-blue-400">
                                    {stat.employee.FirstName[0]}{stat.employee.LastName[0]}
                                </div>
                                <div>
                                    <div className="font-medium text-sm">{stat.employee.FirstName} {stat.employee.LastName}</div>
                                    <div className="text-xs text-gray-500">{stat.employee.Title} · ₺{stat.employee.HourlyRate}/s</div>
                                </div>
                            </div>
                            <div className="text-center text-sm font-medium">{stat.totalShifts}</div>
                            <div className="text-center text-sm">{stat.totalHours.toFixed(0)}s</div>
                            <div className="text-center">
                                {stat.overtimeHours > 0 ? (
                                    <span className="text-sm font-medium text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-md">
                                        +{stat.overtimeHours.toFixed(1)}s
                                    </span>
                                ) : (
                                    <span className="text-sm text-gray-600">—</span>
                                )}
                            </div>
                            <div className="text-right text-sm text-gray-300">₺{stat.normalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                            <div className="text-right text-sm">
                                {stat.overtimeCost > 0 ? (
                                    <span className="text-amber-400 font-medium">₺{stat.overtimeCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</span>
                                ) : (
                                    <span className="text-gray-600">—</span>
                                )}
                            </div>
                            <div className="text-right text-sm font-bold text-white">₺{stat.totalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                        </div>
                    ))}

                    {/* Total Row */}
                    <div className="grid grid-cols-8 gap-2 px-5 py-4 bg-white/[0.03] font-bold text-sm">
                        <div className="col-span-2 text-gray-300">TOPLAM</div>
                        <div className="text-center">{empStats.reduce((s, e) => s + e.totalShifts, 0)}</div>
                        <div className="text-center">{empStats.reduce((s, e) => s + e.totalHours, 0).toFixed(0)}s</div>
                        <div className="text-center text-amber-400">{totalOvertimeHours > 0 ? `+${totalOvertimeHours.toFixed(1)}s` : '—'}</div>
                        <div className="text-right text-gray-300">₺{totalNormalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                        <div className="text-right text-amber-400">₺{totalOvertimeCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                        <div className="text-right text-emerald-400">₺{totalCost.toLocaleString('tr-TR', { maximumFractionDigits: 0 })}</div>
                    </div>
                </div>
            )}

            {/* Info */}
            <div className="glass-card p-4 text-xs text-gray-500 flex items-start gap-2">
                <TrendingUp className="w-4 h-4 flex-shrink-0 mt-0.5" />
                <div>
                    Ek mesai çarpanı: <span className="text-amber-400 font-medium">x{overtimeMultiplier}</span>.
                    8 saati aşan nöbetler ek mesai olarak hesaplanır.
                    Puantaj kaydı yapılmamış nöbetler 8 saat standart olarak kabul edilir.
                </div>
            </div>
        </div>
    );
};

export default OvertimeReport;
