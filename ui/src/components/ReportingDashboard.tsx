import React, { useState, useEffect } from 'react';
import { BarChart3, ChevronLeft, ChevronRight, Loader2, Users, Clock, CalendarOff, TrendingUp } from 'lucide-react';
import { reportApi, departmentApi, employeeApi } from '../services/api';
import type { Department, Employee, TrendData, CostAnalysisReport } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const ReportingDashboard: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [departments, setDepartments] = useState<Department[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [loading, setLoading] = useState(true);

    // Report data
    const [workHoursData, setWorkHoursData] = useState<any>(null);
    const [absenceData, setAbsenceData] = useState<any>(null);
    const [trendData, setTrendData] = useState<TrendData[]>([]);
    const [costAnalysisData, setCostAnalysisData] = useState<CostAnalysisReport | null>(null);

    // Employee summary
    const [selectedEmployee, setSelectedEmployee] = useState<number>(0);
    const [employeeSummary, setEmployeeSummary] = useState<any>(null);

    useEffect(() => {
        Promise.all([departmentApi.list(), employeeApi.list({ page: 1, limit: 1000 })])
            .then(([dRes, eRes]) => {
                setDepartments(dRes.data);
                setEmployees(eRes.data.data);
                if (dRes.data.length > 0 && selectedDept === 0) setSelectedDept(dRes.data[0].ID);
            }).catch(console.error);
    }, []);

    const fetchReports = async () => {
        setLoading(true);
        try {
            const [wh, ab, tr, ca] = await Promise.allSettled([
                reportApi.workHours(month, year, selectedDept || undefined),
                reportApi.absences(month, year, selectedDept || undefined),
                reportApi.trends(1, month, year, selectedDept || undefined),
                reportApi.costAnalysis(month, year, selectedDept || undefined),
            ]);
            if (wh.status === 'fulfilled') setWorkHoursData(wh.value.data);
            if (ab.status === 'fulfilled') setAbsenceData(ab.value.data);
            if (tr.status === 'fulfilled') setTrendData(tr.value.data || []);
            if (ca.status === 'fulfilled') setCostAnalysisData(ca.value.data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (selectedDept) fetchReports();
    }, [month, year, selectedDept]);

    // Employee summary fetch
    useEffect(() => {
        if (selectedEmployee > 0) {
            reportApi.employeeSummary(selectedEmployee, month, year)
                .then(res => setEmployeeSummary(res.data))
                .catch(console.error);
        }
    }, [selectedEmployee, month, year]);

    const prevMonth = () => { if (month === 1) { setMonth(12); setYear(year - 1); } else setMonth(month - 1); };
    const nextMonth = () => { if (month === 12) { setMonth(1); setYear(year + 1); } else setMonth(month + 1); };

    const deptEmployees = employees.filter(e => !selectedDept || e.DepartmentID === selectedDept);

    // Max hours for bar chart scaling
    const maxTrendHours = Math.max(...trendData.map(t => t.total_hours), 1);

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Raporlar & Analiz</h2>
                    <select className="glass-input py-1.5 px-3 text-sm" value={selectedDept} onChange={(e) => setSelectedDept(Number(e.target.value))}>
                        <option value={0} disabled>Bölüm Seçin</option>
                        {departments.map(d => (<option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>))}
                    </select>
                </div>
                <div className="flex items-center gap-2">
                    <button onClick={prevMonth} className="btn-ghost p-2"><ChevronLeft className="w-4 h-4" /></button>
                    <span className="min-w-[160px] text-center font-semibold text-lg">{MONTHS_TR[month - 1]} {year}</span>
                    <button onClick={nextMonth} className="btn-ghost p-2"><ChevronRight className="w-4 h-4" /></button>
                </div>
            </div>

            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" /> Yükleniyor...
                </div>
            ) : (
                <>
                    {/* Summary Cards */}
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                            <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><Clock className="w-3 h-3" /> Toplam Çalışma</div>
                            <div className="text-2xl font-bold text-blue-400">{(workHoursData?.total_hours || 0).toFixed(1)}s</div>
                        </div>
                        <div className="glass-card p-4 bg-gradient-to-br from-emerald-500/5 to-emerald-500/10 border-emerald-500/10">
                            <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><Users className="w-3 h-3" /> İş Günü</div>
                            <div className="text-2xl font-bold text-emerald-400">{workHoursData?.working_days || 0}</div>
                        </div>
                        <div className="glass-card p-4 bg-gradient-to-br from-red-500/5 to-red-500/10 border-red-500/10">
                            <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><CalendarOff className="w-3 h-3" /> İzin Günleri</div>
                            <div className="text-2xl font-bold text-red-400">{absenceData?.total_days || 0}</div>
                        </div>
                        <div className="glass-card p-4 bg-gradient-to-br from-purple-500/5 to-purple-500/10 border-purple-500/10">
                            <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><BarChart3 className="w-3 h-3" /> Personel</div>
                            <div className="text-2xl font-bold text-purple-400">{(workHoursData?.employees?.length || 0)}</div>
                        </div>
                    </div>

                    {/* Trend Chart */}
                    {trendData.length > 0 && (
                        <div className="glass-card p-5">
                            <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center gap-2">
                                <TrendingUp className="w-4 h-4 text-blue-400" /> Aylık Trend ({year})
                            </h3>
                            <div className="flex items-end gap-2 h-40">
                                {trendData.map((t, i) => {
                                    const height = (t.total_hours / maxTrendHours) * 100;
                                    const overtimeHeight = (t.overtime_hours / maxTrendHours) * 100;
                                    return (
                                        <div key={i} className="flex-1 flex flex-col items-center gap-1">
                                            <div className="w-full relative" style={{ height: '100%' }}>
                                                <div
                                                    className="absolute bottom-0 w-full bg-gradient-to-t from-blue-500/30 to-blue-500/10 rounded-t-md transition-all duration-500"
                                                    style={{ height: `${height}%` }}
                                                />
                                                {t.overtime_hours > 0 && (
                                                    <div
                                                        className="absolute bottom-0 w-full bg-gradient-to-t from-amber-500/40 to-amber-500/10 rounded-t-md transition-all duration-500"
                                                        style={{ height: `${overtimeHeight}%` }}
                                                    />
                                                )}
                                            </div>
                                            <div className="text-[10px] text-gray-500">{MONTHS_TR[t.month - 1]?.slice(0, 3)}</div>
                                            <div className="text-[10px] text-gray-400 font-medium">{t.total_hours.toFixed(0)}s</div>
                                        </div>
                                    );
                                })}
                            </div>
                            <div className="flex items-center gap-4 mt-3 text-xs text-gray-500">
                                <span className="flex items-center gap-1"><span className="w-3 h-3 rounded bg-blue-500/30" /> Toplam</span>
                                <span className="flex items-center gap-1"><span className="w-3 h-3 rounded bg-amber-500/40" /> Fazla Mesai</span>
                            </div>
                        </div>
                    )}

                    {/* Employee Summary */}
                    <div className="glass-card p-5 bg-gradient-to-br from-indigo-500/5 to-purple-500/5 border-indigo-500/10">
                        <h3 className="text-sm font-semibold text-gray-300 mb-3">Personel Özeti</h3>
                        <select
                            className="glass-input py-1.5 px-3 text-sm w-full md:w-auto"
                            value={selectedEmployee}
                            onChange={(e) => setSelectedEmployee(Number(e.target.value))}
                        >
                            <option value={0}>Personel seçin...</option>
                            {deptEmployees.map(e => (<option key={e.ID} value={e.ID}>{e.FirstName} {e.LastName}</option>))}
                        </select>
                        {employeeSummary && (
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mt-4">
                                <div className="bg-white/[0.03] rounded-xl p-3 border border-white/5">
                                    <div className="text-xs text-gray-400">Çalışma Saati</div>
                                    <div className="text-lg font-bold text-blue-400 mt-1">{employeeSummary.total_hours?.toFixed(1) || 0}s</div>
                                </div>
                                <div className="bg-white/[0.03] rounded-xl p-3 border border-white/5">
                                    <div className="text-xs text-gray-400">Fazla Mesai</div>
                                    <div className="text-lg font-bold text-amber-400 mt-1">{employeeSummary.overtime_hours?.toFixed(1) || 0}s</div>
                                </div>
                                <div className="bg-white/[0.03] rounded-xl p-3 border border-white/5">
                                    <div className="text-xs text-gray-400">İzin Günleri</div>
                                    <div className="text-lg font-bold text-red-400 mt-1">{employeeSummary.leave_days || 0}</div>
                                </div>
                                <div className="bg-white/[0.03] rounded-xl p-3 border border-white/5">
                                    <div className="text-xs text-gray-400">İş Günleri</div>
                                    <div className="text-lg font-bold text-emerald-400 mt-1">{employeeSummary.working_days || 0}</div>
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Work Hours Table */}
                    {workHoursData?.employees?.length > 0 && (
                        <div className="glass-card overflow-hidden">
                            <div className="px-5 py-3 bg-white/[0.02] border-b border-white/5">
                                <span className="text-sm font-semibold text-gray-300">Çalışma Saatleri Detay</span>
                            </div>
                            <div className="grid grid-cols-5 gap-2 px-5 py-2 bg-white/[0.01] text-xs font-semibold text-gray-400 uppercase">
                                <div className="col-span-2">Personel</div>
                                <div className="text-center">İş Günü</div>
                                <div className="text-center">Toplam Saat</div>
                                <div className="text-center">Ort. Günlük</div>
                            </div>
                            {workHoursData.employees.map((e: any, idx: number) => (
                                <div
                                    key={e.employee_id}
                                    className="grid grid-cols-5 gap-2 px-5 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors items-center animate-slide-up"
                                    style={{ animationDelay: `${idx * 30}ms` }}
                                >
                                    <div className="col-span-2 flex items-center gap-3">
                                        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/15 to-purple-500/15 flex items-center justify-center text-xs font-bold text-blue-400">
                                            {e.employee_name?.split(' ').map((n: string) => n[0]).join('').slice(0, 2)}
                                        </div>
                                        <div>
                                            <div className="font-medium text-sm">{e.employee_name}</div>
                                            <div className="text-xs text-gray-500">{e.department}</div>
                                        </div>
                                    </div>
                                    <div className="text-center text-sm">{e.working_days}</div>
                                    <div className="text-center text-sm font-medium">{e.total_hours?.toFixed(1)}s</div>
                                    <div className="text-center text-sm text-gray-400">{e.avg_daily_hours?.toFixed(1)}s</div>
                                </div>
                            ))}
                        </div>
                    )}
                    {/* Cost Analysis Table */}
                    {costAnalysisData?.employees && costAnalysisData.employees.length > 0 && (
                        <div className="glass-card overflow-hidden mt-6">
                            <div className="px-5 py-3 bg-gradient-to-r from-green-500/10 to-emerald-500/5 border-b border-green-500/20 flex justify-between items-center">
                                <span className="text-sm font-semibold text-green-400">Maliyet Analizi Raporu</span>
                                <span className="text-sm font-bold text-green-300">Toplam: ₺{costAnalysisData.total_cost.toLocaleString('tr-TR')}</span>
                            </div>
                            <div className="grid grid-cols-5 gap-2 px-5 py-2 bg-white/[0.01] text-xs font-semibold text-gray-400 uppercase">
                                <div className="col-span-2">Personel</div>
                                <div className="text-center">Toplam Saat</div>
                                <div className="text-center">Saatlik Ücret</div>
                                <div className="text-right">Aylık Maliyet</div>
                            </div>
                            {costAnalysisData.employees.map((detail: any, idx: number) => (
                                <div
                                    key={detail.employee_id}
                                    className="grid grid-cols-5 gap-2 px-5 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors items-center animate-slide-up"
                                    style={{ animationDelay: `${idx * 30}ms` }}
                                >
                                    <div className="col-span-2 flex items-center gap-3">
                                        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-green-500/15 to-emerald-500/15 flex items-center justify-center text-xs font-bold text-green-400">
                                            {detail.employee_name.split(' ').map((n: string) => n[0]).join('').slice(0, 2)}
                                        </div>
                                        <div>
                                            <div className="font-medium text-sm text-gray-200">{detail.employee_name}</div>
                                            <div className="text-xs text-gray-500">{detail.department}</div>
                                        </div>
                                    </div>
                                    <div className="text-center text-sm">{detail.total_hours.toFixed(1)}s</div>
                                    <div className="text-center text-sm text-gray-400">₺{detail.hourly_rate}</div>
                                    <div className="text-right text-sm font-medium text-green-400">₺{detail.total_cost.toLocaleString('tr-TR')}</div>
                                </div>
                            ))}
                        </div>
                    )}
                </>
            )}
        </div>
    );
};

export default ReportingDashboard;
