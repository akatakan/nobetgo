import React, { useState, useEffect } from 'react';
import { TrendingUp, ChevronLeft, ChevronRight, Loader2, DollarSign, Clock, Users, Settings } from 'lucide-react';
import { overtimeApi, overtimeRuleApi, publicHolidayApi, departmentApi } from '../services/api';
import type { OvertimeSummary, OvertimeRule, PublicHoliday, Department } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const OvertimeReport: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [departments, setDepartments] = useState<Department[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [summaries, setSummaries] = useState<OvertimeSummary[]>([]);
    const [loading, setLoading] = useState(true);

    // Rules & Holidays
    const [showRules, setShowRules] = useState(false);
    const [rules, setRules] = useState<OvertimeRule[]>([]);
    const [holidays, setHolidays] = useState<PublicHoliday[]>([]);
    const [ruleForm, setRuleForm] = useState({
        name: '', weekly_hour_limit: 45, daily_hour_limit: 11,
        overtime_multiplier: 1.5, weekend_multiplier: 2.0,
        holiday_multiplier: 2.5, night_shift_extra: 0.1, is_active: true,
    });
    const [holidayForm, setHolidayForm] = useState({ name: '', date: '' });

    useEffect(() => {
        departmentApi.list().then(res => {
            setDepartments(res.data);
            if (res.data.length > 0 && selectedDept === 0) setSelectedDept(res.data[0].ID);
        }).catch(console.error);
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await overtimeApi.departmentSummary(selectedDept, month, year);
            setSummaries(res.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (selectedDept) fetchData();
    }, [month, year, selectedDept]);

    const fetchRulesAndHolidays = async () => {
        try {
            const [rRes, hRes] = await Promise.all([
                overtimeRuleApi.list(),
                publicHolidayApi.list(year),
            ]);
            setRules(rRes.data || []);
            setHolidays(hRes.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        if (showRules) fetchRulesAndHolidays();
    }, [showRules, year]);

    const prevMonth = () => { if (month === 1) { setMonth(12); setYear(year - 1); } else setMonth(month - 1); };
    const nextMonth = () => { if (month === 12) { setMonth(1); setYear(year + 1); } else setMonth(month + 1); };

    const totalNormal = summaries.reduce((s, e) => s + e.normal_hours, 0);
    const totalOvertime = summaries.reduce((s, e) => s + e.overtime_hours, 0);
    const totalWeekend = summaries.reduce((s, e) => s + e.weekend_hours, 0);
    const totalHoliday = summaries.reduce((s, e) => s + e.holiday_hours, 0);

    const handleCreateRule = async () => {
        if (!ruleForm.name) return;
        await overtimeRuleApi.create(ruleForm);
        setRuleForm({ name: '', weekly_hour_limit: 45, daily_hour_limit: 11, overtime_multiplier: 1.5, weekend_multiplier: 2.0, holiday_multiplier: 2.5, night_shift_extra: 0.1, is_active: true });
        fetchRulesAndHolidays();
    };

    const handleCreateHoliday = async () => {
        if (!holidayForm.name || !holidayForm.date) return;
        await publicHolidayApi.create(holidayForm);
        setHolidayForm({ name: '', date: '' });
        fetchRulesAndHolidays();
    };

    const handleDeleteRule = async (id: number) => {
        await overtimeRuleApi.delete(id);
        fetchRulesAndHolidays();
    };

    const handleDeleteHoliday = async (id: number) => {
        await publicHolidayApi.delete(id);
        fetchRulesAndHolidays();
    };

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Mesai & Fazla Çalışma</h2>
                    <select className="glass-input py-1.5 px-3 text-sm" value={selectedDept} onChange={(e) => setSelectedDept(Number(e.target.value))}>
                        <option value={0} disabled>Bölüm Seçin</option>
                        {departments.map((d) => (<option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>))}
                    </select>
                </div>
                <div className="flex items-center gap-2">
                    <button onClick={() => setShowRules(!showRules)} className="btn-ghost text-xs py-2 px-3 flex items-center gap-2">
                        <Settings className="w-3.5 h-3.5" /> Kurallar & Tatiller
                    </button>
                    <div className="h-6 w-px bg-white/10 mx-1"></div>
                    <button onClick={prevMonth} className="btn-ghost p-2"><ChevronLeft className="w-4 h-4" /></button>
                    <span className="min-w-[160px] text-center font-semibold text-lg">{MONTHS_TR[month - 1]} {year}</span>
                    <button onClick={nextMonth} className="btn-ghost p-2"><ChevronRight className="w-4 h-4" /></button>
                </div>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><Clock className="w-3 h-3" /> Normal Saat</div>
                    <div className="text-2xl font-bold text-blue-400">{totalNormal.toFixed(1)}s</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-amber-500/5 to-amber-500/10 border-amber-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><TrendingUp className="w-3 h-3" /> Fazla Mesai</div>
                    <div className="text-2xl font-bold text-amber-400">{totalOvertime.toFixed(1)}s</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-purple-500/5 to-purple-500/10 border-purple-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><DollarSign className="w-3 h-3" /> Hafta Sonu</div>
                    <div className="text-2xl font-bold text-purple-400">{totalWeekend.toFixed(1)}s</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-red-500/5 to-red-500/10 border-red-500/10">
                    <div className="text-xs text-gray-400 mb-1 flex items-center gap-1"><Users className="w-3 h-3" /> Tatil Mesaisi</div>
                    <div className="text-2xl font-bold text-red-400">{totalHoliday.toFixed(1)}s</div>
                </div>
            </div>

            {/* Rules & Holidays Panel */}
            {showRules && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5 animate-slide-down">
                    {/* Rules */}
                    <div className="glass-card p-5">
                        <h3 className="text-sm font-semibold text-gray-300 mb-3">Mesai Kuralları</h3>
                        {rules.map(r => (
                            <div key={r.ID} className="flex items-center justify-between py-2 border-b border-white/5 last:border-0">
                                <div>
                                    <div className="text-sm font-medium">{r.name}</div>
                                    <div className="text-xs text-gray-500">
                                        Haftalık: {r.weekly_hour_limit}s · Günlük: {r.daily_hour_limit}s · Çarpan: x{r.overtime_multiplier}
                                    </div>
                                </div>
                                <button onClick={() => handleDeleteRule(r.ID)} className="btn-ghost text-xs text-red-400/60 hover:text-red-400 p-1">Sil</button>
                            </div>
                        ))}
                        <div className="mt-3 space-y-2">
                            <input className="glass-input w-full text-sm" placeholder="Kural adı" value={ruleForm.name} onChange={e => setRuleForm({ ...ruleForm, name: e.target.value })} />
                            <div className="grid grid-cols-2 gap-2">
                                <input type="number" className="glass-input text-sm" placeholder="Haftalık limit" value={ruleForm.weekly_hour_limit} onChange={e => setRuleForm({ ...ruleForm, weekly_hour_limit: Number(e.target.value) })} />
                                <input type="number" className="glass-input text-sm" placeholder="Çarpan" step="0.1" value={ruleForm.overtime_multiplier} onChange={e => setRuleForm({ ...ruleForm, overtime_multiplier: Number(e.target.value) })} />
                            </div>
                            <button onClick={handleCreateRule} className="btn-primary text-xs w-full">Kural Ekle</button>
                        </div>
                    </div>
                    {/* Holidays */}
                    <div className="glass-card p-5">
                        <h3 className="text-sm font-semibold text-gray-300 mb-3">Resmi Tatiller ({year})</h3>
                        {holidays.map(h => (
                            <div key={h.ID} className="flex items-center justify-between py-2 border-b border-white/5 last:border-0">
                                <div>
                                    <div className="text-sm font-medium">{h.name}</div>
                                    <div className="text-xs text-gray-500">{h.date?.split('T')[0]}</div>
                                </div>
                                <button onClick={() => handleDeleteHoliday(h.ID)} className="btn-ghost text-xs text-red-400/60 hover:text-red-400 p-1">Sil</button>
                            </div>
                        ))}
                        <div className="mt-3 space-y-2">
                            <input className="glass-input w-full text-sm" placeholder="Tatil adı" value={holidayForm.name} onChange={e => setHolidayForm({ ...holidayForm, name: e.target.value })} />
                            <input type="date" className="glass-input w-full text-sm" value={holidayForm.date} onChange={e => setHolidayForm({ ...holidayForm, date: e.target.value })} />
                            <button onClick={handleCreateHoliday} className="btn-primary text-xs w-full">Tatil Ekle</button>
                        </div>
                    </div>
                </div>
            )}

            {/* Employee Table */}
            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />Yükleniyor...
                </div>
            ) : summaries.length === 0 ? (
                <div className="glass-card p-12 text-center text-gray-500">
                    <Users className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Bu dönem için veri bulunamadı.</p>
                </div>
            ) : (
                <div className="glass-card overflow-hidden">
                    <div className="grid grid-cols-8 gap-2 px-5 py-3 bg-white/[0.02] border-b border-white/5 text-xs font-semibold text-gray-400 uppercase">
                        <div className="col-span-2">Personel</div>
                        <div className="text-center">İş Günü</div>
                        <div className="text-center">Normal</div>
                        <div className="text-center">Fazla Mesai</div>
                        <div className="text-center">Hafta Sonu</div>
                        <div className="text-center">Tatil</div>
                        <div className="text-right">Toplam</div>
                    </div>
                    {summaries.map((s, idx) => (
                        <div
                            key={s.employee_id}
                            className="grid grid-cols-8 gap-2 px-5 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors items-center animate-slide-up"
                            style={{ animationDelay: `${idx * 30}ms`, animationFillMode: 'both' }}
                        >
                            <div className="col-span-2 flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/15 to-purple-500/15 flex items-center justify-center text-xs font-bold text-blue-400">
                                    {s.employee_name?.split(' ').map(n => n[0]).join('').slice(0, 2)}
                                </div>
                                <div className="font-medium text-sm">{s.employee_name}</div>
                            </div>
                            <div className="text-center text-sm">{s.working_days}</div>
                            <div className="text-center text-sm text-gray-300">{s.normal_hours.toFixed(1)}s</div>
                            <div className="text-center">
                                {s.overtime_hours > 0 ? (
                                    <span className="text-sm font-medium text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-md">+{s.overtime_hours.toFixed(1)}s</span>
                                ) : <span className="text-sm text-gray-600">—</span>}
                            </div>
                            <div className="text-center">
                                {s.weekend_hours > 0 ? (
                                    <span className="text-sm font-medium text-purple-400">{s.weekend_hours.toFixed(1)}s</span>
                                ) : <span className="text-sm text-gray-600">—</span>}
                            </div>
                            <div className="text-center">
                                {s.holiday_hours > 0 ? (
                                    <span className="text-sm font-medium text-red-400">{s.holiday_hours.toFixed(1)}s</span>
                                ) : <span className="text-sm text-gray-600">—</span>}
                            </div>
                            <div className="text-right text-sm font-bold text-white">{s.total_hours.toFixed(1)}s</div>
                        </div>
                    ))}
                    {/* Total Row */}
                    <div className="grid grid-cols-8 gap-2 px-5 py-4 bg-white/[0.03] font-bold text-sm">
                        <div className="col-span-2 text-gray-300">TOPLAM</div>
                        <div className="text-center">—</div>
                        <div className="text-center text-gray-300">{totalNormal.toFixed(1)}s</div>
                        <div className="text-center text-amber-400">{totalOvertime > 0 ? `+${totalOvertime.toFixed(1)}s` : '—'}</div>
                        <div className="text-center text-purple-400">{totalWeekend > 0 ? `${totalWeekend.toFixed(1)}s` : '—'}</div>
                        <div className="text-center text-red-400">{totalHoliday > 0 ? `${totalHoliday.toFixed(1)}s` : '—'}</div>
                        <div className="text-right text-emerald-400">{(totalNormal + totalOvertime + totalWeekend + totalHoliday).toFixed(1)}s</div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OvertimeReport;
