import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Calendar, Loader2 } from 'lucide-react';
import { scheduleApi, departmentApi } from '../services/api';
import type { Schedule, Department } from '../types';

const DAYS_TR = ['Pzt', 'Sal', 'Çar', 'Per', 'Cum', 'Cmt', 'Paz'];
const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const ScheduleViewer: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [data, setData] = useState<Schedule[]>([]);
    const [loading, setLoading] = useState(true);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);

    useEffect(() => {
        departmentApi.list().then(res => {
            setDepartments(res.data);
            if (res.data.length > 0 && selectedDept === 0) {
                setSelectedDept(res.data[0].ID);
            }
        }).catch(console.error);
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await scheduleApi.getMonthly(month, year, selectedDept || undefined);
            setData(res.data || []);
        } catch (err) {
            console.error(err);
            setData([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [month, year, selectedDept]);

    const prevMonth = () => {
        if (month === 1) { setMonth(12); setYear(year - 1); }
        else setMonth(month - 1);
    };
    const nextMonth = () => {
        if (month === 12) { setMonth(1); setYear(year + 1); }
        else setMonth(month + 1);
    };

    // Build calendar grid
    const firstDay = new Date(year, month - 1, 1);
    const daysInMonth = new Date(year, month, 0).getDate();
    // Monday-indexed: 0=Mon ... 6=Sun
    let startOffset = firstDay.getDay() - 1;
    if (startOffset < 0) startOffset = 6;

    const calendarCells: (number | null)[] = [];
    for (let i = 0; i < startOffset; i++) calendarCells.push(null);
    for (let d = 1; d <= daysInMonth; d++) calendarCells.push(d);
    while (calendarCells.length % 7 !== 0) calendarCells.push(null);

    // Group schedules by day
    const schedulesByDay: Record<number, Schedule[]> = {};
    data.forEach((s) => {
        const day = new Date(s.Date).getDate();
        if (!schedulesByDay[day]) schedulesByDay[day] = [];
        schedulesByDay[day].push(s);
    });

    const today = new Date();
    const isToday = (day: number) =>
        day === today.getDate() && month === today.getMonth() + 1 && year === today.getFullYear();

    const currentDept = departments.find(d => d.ID === selectedDept);

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">Nöbet Takvimi</h2>
                    <select
                        className="glass-input py-1.5 px-3 text-sm"
                        value={selectedDept}
                        onChange={(e) => setSelectedDept(Number(e.target.value))}
                    >
                        <option value={0}>Tüm Bölümler</option>
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

            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />
                    Takvim yükleniyor...
                </div>
            ) : (
                <div className="glass-card overflow-hidden">
                    {/* Day headers */}
                    <div className="grid grid-cols-7 border-b border-white/5">
                        {DAYS_TR.map((d) => (
                            <div key={d} className="px-2 py-3 text-center text-xs font-semibold text-gray-400 uppercase tracking-wider">
                                {d}
                            </div>
                        ))}
                    </div>

                    {/* Calendar grid */}
                    <div className="grid grid-cols-7">
                        {calendarCells.map((day, idx) => {
                            const schedules = day ? schedulesByDay[day] || [] : [];
                            const weekend = idx % 7 >= 5; // Sat, Sun
                            return (
                                <div
                                    key={idx}
                                    className={`min-h-[100px] p-2 border-b border-r border-white/[0.03] transition-colors
                                        ${day ? 'hover:bg-white/[0.02]' : 'bg-white/[0.01]'}
                                        ${weekend && day ? 'bg-blue-500/[0.02]' : ''}
                                        ${isToday(day || 0) ? 'bg-blue-500/[0.06] ring-1 ring-inset ring-blue-500/20' : ''}
                                    `}
                                >
                                    {day && (
                                        <>
                                            <div className={`text-sm font-medium mb-1 ${isToday(day) ? 'text-blue-400' : weekend ? 'text-gray-400' : 'text-gray-300'}`}>
                                                {day}
                                            </div>
                                            <div className="space-y-1">
                                                {schedules.slice(0, 3).map((s) => (
                                                    <div
                                                        key={s.ID}
                                                        className="text-[10px] px-1.5 py-0.5 rounded-md truncate cursor-default transition-opacity hover:opacity-90"
                                                        style={{
                                                            backgroundColor: (s.ShiftType?.Color || '#3b82f6') + '25',
                                                            color: s.ShiftType?.Color || '#3b82f6',
                                                            borderLeft: `2px solid ${s.ShiftType?.Color || '#3b82f6'}`,
                                                        }}
                                                        title={`${s.Employee?.FirstName} ${s.Employee?.LastName} — ${s.ShiftType?.Name}`}
                                                    >
                                                        {s.Employee?.FirstName?.[0]}.{s.Employee?.LastName?.[0]}. {s.ShiftType?.Name}
                                                    </div>
                                                ))}
                                                {schedules.length > 3 && (
                                                    <div className="text-[10px] text-gray-500 pl-1">
                                                        +{schedules.length - 3} daha
                                                    </div>
                                                )}
                                            </div>
                                        </>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>
            )}

            {/* Stats */}
            {!loading && (
                <div className="flex gap-4 text-sm text-gray-400 animate-slide-up">
                    {currentDept && (
                        <div className="glass-card px-4 py-2">
                            Bölüm: <span className="text-blue-400 font-semibold">{currentDept.Floor}. Kat - {currentDept.Name}</span>
                        </div>
                    )}
                    <div className="glass-card px-4 py-2">
                        <Calendar className="w-4 h-4 inline mr-2" />
                        Toplam: <span className="text-white font-semibold">{data.length}</span> nöbet
                    </div>
                    <div className="glass-card px-4 py-2">
                        Personel: <span className="text-white font-semibold">
                            {new Set(data.map(s => s.EmployeeID)).size}
                        </span> kişi
                    </div>
                </div>
            )}
        </div>
    );
};

export default ScheduleViewer;
