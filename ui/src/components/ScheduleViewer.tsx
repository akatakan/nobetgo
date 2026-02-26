import React, { useState, useEffect, useRef } from 'react';
import { ChevronLeft, ChevronRight, Calendar, Loader2, Trash2, Printer, FileSpreadsheet, AlertCircle } from 'lucide-react';
import { scheduleApi, departmentApi } from '../services/api';
import type { Schedule, Department } from '../types';
import * as XLSX from 'xlsx';
import ScheduleEditModal from './ScheduleEditModal';

const DAYS_TR = ['Pzt', 'Sal', 'Çar', 'Per', 'Cum', 'Cmt', 'Paz'];
const DAYS_FULL_TR = ['Pazartesi', 'Salı', 'Çarşamba', 'Perşembe', 'Cuma', 'Cumartesi', 'Pazar'];
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
    const [clearing, setClearing] = useState(false);
    const [showClearConfirm, setShowClearConfirm] = useState(false);
    const [draggedSchedule, setDraggedSchedule] = useState<Schedule | null>(null);
    const [modalState, setModalState] = useState<{ isOpen: boolean, schedule: Schedule | null, date: Date | null }>({
        isOpen: false, schedule: null, date: null
    });
    const printRef = useRef<HTMLDivElement>(null);

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

    const handleClear = async () => {
        setClearing(true);
        try {
            await scheduleApi.clear(month, year);
            setShowClearConfirm(false);
            fetchData();
        } catch (err) {
            console.error(err);
            alert('Silme hatası: ' + err);
        } finally {
            setClearing(false);
        }
    };

    // Drag and Drop Handlers
    const handleDragStart = (e: React.DragEvent, schedule: Schedule) => {
        setDraggedSchedule(schedule);
        e.dataTransfer.effectAllowed = 'move';
        // semi-transparent ghost image
        if (e.target instanceof HTMLElement) e.target.style.opacity = '0.5';
    };

    const handleDragEnd = (e: React.DragEvent) => {
        if (e.target instanceof HTMLElement) e.target.style.opacity = '1';
    };

    const handleDragOver = (e: React.DragEvent) => {
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
    };

    const handleDrop = async (e: React.DragEvent, targetDay: number) => {
        e.preventDefault();
        if (!draggedSchedule) return;

        // If dropped on the same day, do nothing
        const currentDay = new Date(draggedSchedule.Date).getDate();
        if (currentDay === targetDay) return;

        // Calculate new date
        // Create date object for the target day at noon to avoid timezone shift issues
        const targetDate = new Date(year, month - 1, targetDay, 12, 0, 0);
        const targetDateStr = targetDate.toISOString();

        // Optimistic update
        const updated = { ...draggedSchedule, Date: targetDateStr };
        setData(prev => prev.map(s => s.ID === draggedSchedule.ID ? updated : s));

        try {
            // Need to send the full object with updated date
            await scheduleApi.update(draggedSchedule.ID, { ...draggedSchedule, Date: targetDateStr });
            // No need to fetch immediately if optimistic update worked, but good for consistency
        } catch (err) {
            console.error(err);
            alert('Taşıma başarısız!');
            fetchData(); // Revert on error
        } finally {
            setDraggedSchedule(null);
        }
    };

    // Print functionality
    const handlePrint = () => {
        const deptName = currentDept ? `${currentDept.Floor}. Kat - ${currentDept.Name}` : 'Tüm Bölümler';
        const title = `Nöbet Çizelgesi — ${MONTHS_TR[month - 1]} ${year} — ${deptName}`;

        // Build print-friendly HTML table
        let html = `
        <html><head><title>${title}</title>
        <style>
            * { margin: 0; padding: 0; box-sizing: border-box; }
            body { font-family: 'Segoe UI', Arial, sans-serif; padding: 20px; color: #111; }
            h1 { font-size: 18px; text-align: center; margin-bottom: 4px; }
            h2 { font-size: 13px; text-align: center; color: #666; margin-bottom: 16px; font-weight: normal; }
            table { width: 100%; border-collapse: collapse; table-layout: fixed; }
            th { background: #f0f0f0; padding: 8px 4px; border: 1px solid #ccc; font-size: 11px; text-align: center; }
            td { border: 1px solid #ccc; padding: 4px; vertical-align: top; min-height: 70px; height: 70px; font-size: 10px; }
            .day-num { font-weight: bold; font-size: 12px; margin-bottom: 3px; }
            .shift-entry { padding: 2px 4px; margin: 1px 0; border-radius: 3px; background: #f5f5f5; border-left: 3px solid #3b82f6; font-size: 9px; }
            .shift-name { font-weight: 600; }
            .shift-time { color: #666; font-size: 8px; }
            .weekend { background: #fafafa; }
            .empty { background: #f9f9f9; }
            .stats { margin-top: 16px; font-size: 11px; color: #555; }
            .stats table { width: auto; margin-top: 8px; }
            .stats td, .stats th { padding: 4px 12px; font-size: 11px; }
            @media print { body { padding: 10px; } }
        </style></head><body>
        <h1>${title}</h1>
        <h2>Oluşturulma: ${new Date().toLocaleDateString('tr-TR')} — Toplam ${data.length} atama</h2>
        <table><thead><tr>`;

        DAYS_FULL_TR.forEach(d => { html += `<th>${d}</th>`; });
        html += '</tr></thead><tbody>';

        for (let row = 0; row < calendarCells.length; row += 7) {
            html += '<tr>';
            for (let col = 0; col < 7; col++) {
                const day = calendarCells[row + col];
                const isWeekend = col >= 5;
                if (day) {
                    const schedules = schedulesByDay[day] || [];
                    html += `<td class="${isWeekend ? 'weekend' : ''}">`;
                    html += `<div class="day-num">${day}</div>`;
                    schedules.forEach(s => {
                        const name = s.Employee ? `${s.Employee.FirstName} ${s.Employee.LastName}` : '!!! BOŞ NÖBET !!!';
                        const shift = s.ShiftType?.Name || '';
                        const time = s.ShiftType ? `${s.ShiftType.StartTime}–${s.ShiftType.EndTime}` : '';
                        html += `<div class="shift-entry" style="${!s.Employee ? 'border-left-color:red; background:#fff1f1;' : ''}"><span class="shift-name">${name}</span><br/><span class="shift-time">${shift} ${time}</span></div>`;
                    });
                    html += '</td>';
                } else {
                    html += '<td class="empty"></td>';
                }
            }
            html += '</tr>';
        }
        html += '</tbody></table>';

        // Stats section — shifts per employee
        const empShiftCount: Record<string, number> = {};
        data.forEach(s => {
            const name = `${s.Employee?.FirstName} ${s.Employee?.LastName}`;
            empShiftCount[name] = (empShiftCount[name] || 0) + 1;
        });

        html += '<div class="stats"><strong>Personel Nöbet Dağılımı:</strong><table><thead><tr><th>Personel</th><th>Nöbet Sayısı</th></tr></thead><tbody>';
        Object.entries(empShiftCount).sort((a, b) => b[1] - a[1]).forEach(([name, count]) => {
            html += `<tr><td>${name}</td><td style="text-align:center">${count}</td></tr>`;
        });
        html += '</tbody></table></div></body></html>';

        const printWindow = window.open('', '_blank');
        if (printWindow) {
            printWindow.document.write(html);
            printWindow.document.close();
            printWindow.focus();
            setTimeout(() => printWindow.print(), 300);
        }
    };

    // Excel export — calendar grid format
    const handleExcelExport = () => {
        const deptName = currentDept ? `${currentDept.Floor}. Kat - ${currentDept.Name}` : 'Tüm Bölümler';

        // Build calendar grid as 2D array
        const gridData: (string | null)[][] = [];

        // Title row
        gridData.push([`Nöbet Çizelgesi — ${MONTHS_TR[month - 1]} ${year} — ${deptName}`]);
        gridData.push([]); // empty spacer row

        // Header row: weekday names
        gridData.push(DAYS_FULL_TR.map(d => d));

        // Calendar weeks
        for (let row = 0; row < calendarCells.length; row += 7) {
            const weekRow: string[] = [];
            for (let col = 0; col < 7; col++) {
                const day = calendarCells[row + col];
                if (day) {
                    const schedules = schedulesByDay[day] || [];
                    let cellText = `${day}`;
                    schedules.forEach(s => {
                        const name = s.Employee ? `${s.Employee.FirstName} ${s.Employee.LastName}` : 'BOŞ NÖBET';
                        const shift = s.ShiftType?.Name || '';
                        const time = s.ShiftType ? `${s.ShiftType.StartTime}-${s.ShiftType.EndTime}` : '';
                        cellText += `\n${name}\n${shift} ${time}`;
                    });
                    weekRow.push(cellText);
                } else {
                    weekRow.push('');
                }
            }
            gridData.push(weekRow);
        }

        // Empty row before summary
        gridData.push([]);
        gridData.push([]);

        // Summary section
        gridData.push(['Personel Nöbet Dağılımı']);
        gridData.push(['Personel', 'Nöbet Sayısı']);

        const empShiftCount: Record<string, number> = {};
        data.forEach(s => {
            const name = `${s.Employee?.FirstName} ${s.Employee?.LastName}`;
            empShiftCount[name] = (empShiftCount[name] || 0) + 1;
        });
        Object.entries(empShiftCount)
            .sort((a, b) => b[1] - a[1])
            .forEach(([name, count]) => {
                gridData.push([name, String(count)]);
            });

        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.aoa_to_sheet(gridData);

        // Style: merge title row across 7 columns
        ws['!merges'] = [{ s: { r: 0, c: 0 }, e: { r: 0, c: 6 } }];

        // Set column widths (each day column ~22 chars wide)
        ws['!cols'] = Array(7).fill({ wch: 24 });

        // Set row heights for calendar rows to accommodate multi-line content
        const rowHeights: Record<number, { hpt: number }> = {};
        for (let i = 3; i < 3 + Math.ceil(calendarCells.length / 7); i++) {
            rowHeights[i] = { hpt: 60 }; // ~60pt row height for calendar cells
        }
        ws['!rows'] = [];
        for (let i = 0; i < gridData.length; i++) {
            ws['!rows'][i] = rowHeights[i] || { hpt: 16 };
        }
        ws['!rows'][0] = { hpt: 28 }; // Title row taller
        ws['!rows'][2] = { hpt: 22 }; // Header row

        XLSX.utils.book_append_sheet(wb, ws, 'Nöbet Takvimi');

        const fileName = `Nobet_${MONTHS_TR[month - 1]}_${year}_${deptName.replace(/\s/g, '_')}.xlsx`;
        XLSX.writeFile(wb, fileName);
    };

    const handleEditSchedule = (schedule: Schedule) => {
        setModalState({ isOpen: true, schedule, date: new Date(schedule.Date) });
    };

    const handleNewSchedule = (day: number) => {
        const date = new Date(year, month - 1, day);
        setModalState({ isOpen: true, schedule: null, date });
    };

    return (
        <div className="space-y-6 animate-fade-in" ref={printRef}>
            <ScheduleEditModal
                isOpen={modalState.isOpen}
                schedule={modalState.schedule}
                date={modalState.date || new Date()}
                onClose={() => setModalState({ ...modalState, isOpen: false })}
                onSuccess={fetchData}
            />
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4 flex-wrap">
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
                    {data.length > 0 && (
                        <>
                            <button
                                onClick={handlePrint}
                                className="btn-ghost text-blue-400 hover:text-blue-300 hover:bg-blue-500/10 px-3 py-1.5 text-sm"
                                title="Yazdır"
                            >
                                <Printer className="w-4 h-4 inline mr-1" />
                                Yazdır
                            </button>
                            <button
                                onClick={handleExcelExport}
                                className="btn-ghost text-green-400 hover:text-green-300 hover:bg-green-500/10 px-3 py-1.5 text-sm"
                                title="Excel'e Aktar"
                            >
                                <FileSpreadsheet className="w-4 h-4 inline mr-1" />
                                Excel
                            </button>
                            <button
                                onClick={() => setShowClearConfirm(true)}
                                className="btn-ghost text-red-400 hover:text-red-300 hover:bg-red-500/10 px-3 py-1.5 text-sm"
                                title="Bu ayın listesini sil"
                            >
                                <Trash2 className="w-4 h-4 inline mr-1" />
                                Listeyi Sil
                            </button>
                        </>
                    )}
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

            {/* Clear confirmation dialog */}
            {showClearConfirm && (
                <div className="glass-card p-6 mb-6 border border-red-500/20 animate-slide-up">
                    <h4 className="text-lg font-semibold text-red-400 mb-2">Listeyi Silmek İstediğinize Emin Misiniz?</h4>
                    <p className="text-gray-400 text-sm mb-4">
                        {MONTHS_TR[month - 1]} {year} dönemi için oluşturulan <strong className="text-white">{data.length}</strong> nöbet ataması kalıcı olarak silinecek.
                    </p>
                    <div className="flex gap-3">
                        <button
                            onClick={() => setShowClearConfirm(false)}
                            className="btn-ghost text-sm"
                        >
                            Vazgeç
                        </button>
                        <button
                            onClick={handleClear}
                            disabled={clearing}
                            className="px-4 py-2 rounded-lg bg-red-500/20 text-red-400 hover:bg-red-500/30 border border-red-500/30 text-sm transition-colors"
                        >
                            {clearing ? <Loader2 className="w-4 h-4 animate-spin inline mr-2" /> : <Trash2 className="w-4 h-4 inline mr-2" />}
                            {clearing ? 'Siliniyor...' : 'Evet, Sil'}
                        </button>
                    </div>
                </div>
            )}

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
                            const weekend = idx % 7 >= 5;
                            return (
                                <div
                                    key={idx}
                                    onDragOver={handleDragOver}
                                    onDrop={(e) => day && handleDrop(e, day)}
                                    className={`min-h-[100px] p-2 border-b border-r border-white/[0.03] transition-colors group
                                        ${day ? 'hover:bg-white/[0.02]' : 'bg-white/[0.01]'}
                                        ${weekend && day ? 'bg-blue-500/[0.02]' : ''}
                                        ${isToday(day || 0) ? 'bg-blue-500/[0.06] ring-1 ring-inset ring-blue-500/20' : ''}
                                        ${draggedSchedule && day ? 'hover:ring-2 hover:ring-blue-500/30' : ''}
                                    `}
                                >
                                    {day && (
                                        <>
                                            <div className={`text-sm font-medium mb-1 flex justify-between items-center ${isToday(day) ? 'text-blue-400' : weekend ? 'text-gray-400' : 'text-gray-300'}`}>
                                                <span>{day}</span>
                                                <button
                                                    onClick={() => handleNewSchedule(day)}
                                                    className="opacity-0 group-hover:opacity-100 text-[10px] bg-white/10 hover:bg-white/20 rounded px-1.5 py-0.5 transition-all text-gray-300"
                                                    title="Yeni Nöbet Ekle"
                                                >
                                                    +
                                                </button>
                                            </div>
                                            <div className="space-y-1">
                                                {schedules.map((s) => {
                                                    const isUnfilled = !s.EmployeeID || s.EmployeeID === 0 || !s.Employee;
                                                    const fScore = s.Employee?.FatigueScore || 0;
                                                    const isExtremelyTired = !isUnfilled && fScore >= 50;
                                                    const isTired = !isUnfilled && fScore >= 40 && fScore < 50;

                                                    // Heatmap bg overrides shift type color if tired/unfilled
                                                    let bgStatusColor = (s.ShiftType?.Color || '#3b82f6') + '25';
                                                    let borderStatusColor = s.ShiftType?.Color || '#3b82f6';

                                                    if (isUnfilled) {
                                                        bgStatusColor = 'rgba(239, 68, 68, 0.15)'; // Soft red
                                                        borderStatusColor = '#ef4444'; // Solid red
                                                    } else if (isExtremelyTired) {
                                                        bgStatusColor = 'rgba(239, 68, 68, 0.2)';
                                                        borderStatusColor = '#ef4444';
                                                    } else if (isTired) {
                                                        bgStatusColor = 'rgba(245, 158, 11, 0.2)';
                                                        borderStatusColor = '#f59e0b';
                                                    }

                                                    return (
                                                        <div
                                                            key={s.ID}
                                                            draggable
                                                            onDragStart={(e) => handleDragStart(e, s)}
                                                            onDragEnd={handleDragEnd}
                                                            onClick={() => handleEditSchedule(s)}
                                                            className={`text-[10px] px-1.5 py-0.5 rounded-md truncate cursor-pointer hover:scale-[1.02] hover:shadow-md active:cursor-grabbing transition-all group/item relative flex items-center justify-between
                                                                ${isUnfilled ? 'border-dashed border-[1px] animate-pulse' : ''}`}
                                                            style={{
                                                                backgroundColor: bgStatusColor,
                                                                color: borderStatusColor,
                                                                borderLeft: isUnfilled ? undefined : `2px solid ${borderStatusColor}`,
                                                                borderColor: isUnfilled ? borderStatusColor : undefined
                                                            }}
                                                            title={isUnfilled ? `ATANMAMIŞ NÖBET: ${s.ShiftType?.Name}` : `${s.Employee?.FirstName} ${s.Employee?.LastName} — ${s.ShiftType?.Name}`}
                                                        >
                                                            <span className="truncate flex items-center gap-1">
                                                                {isUnfilled ? (
                                                                    <>
                                                                        <AlertCircle className="w-2.5 h-2.5" />
                                                                        <span className="font-bold">BOŞ NÖBET</span>
                                                                    </>
                                                                ) : (
                                                                    `${s.Employee?.FirstName} ${s.Employee?.LastName}`
                                                                )}
                                                            </span>
                                                            {(isTired || isExtremelyTired) && (
                                                                <div
                                                                    className="group/tooltip relative ml-1"
                                                                    aria-label="Tükenmişlik Riski"
                                                                    title={isExtremelyTired ? `Kritik Uyarı: Bu personel sınırın ötesinde nöbet tutuyor (Yorgunluk: ${fScore}). Hata yapma riski yüksek!` : `Dikkat: Bu personel oldukça yorgun (Yorgunluk: ${fScore}).`}
                                                                >
                                                                    <AlertCircle className={`w-3 h-3 ${isExtremelyTired ? 'text-red-500 animate-pulse' : 'text-amber-500'}`} />
                                                                </div>
                                                            )}
                                                        </div>
                                                    )
                                                })}
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
                <div className="flex gap-4 text-sm text-gray-400 animate-slide-up flex-wrap">
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
