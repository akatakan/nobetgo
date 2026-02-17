import React, { useState, useEffect } from 'react';
import { Wand2, CheckCircle2, ChevronRight, Loader2, Users, AlertCircle } from 'lucide-react';
import { scheduleApi, departmentApi, shiftTypeApi, employeeApi } from '../services/api';
import type { ScheduleRequest, Department, ShiftType, Employee } from '../types';

interface Props {
    onNavigate?: (tab: string) => void;
}

const ScheduleWizard: React.FC<Props> = ({ onNavigate }) => {
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [shiftTypes, setShiftTypes] = useState<ShiftType[]>([]);
    const [selectedShiftTypes, setSelectedShiftTypes] = useState<number[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [selectedEmployees, setSelectedEmployees] = useState<number[]>([]);
    const [loadingEmployees, setLoadingEmployees] = useState(false);
    const [params, setParams] = useState<ScheduleRequest>({
        month: new Date().getMonth() + 2 > 12 ? 1 : new Date().getMonth() + 2,
        year: new Date().getMonth() + 2 > 12 ? new Date().getFullYear() + 1 : new Date().getFullYear(),
        department_id: 0,
        shift_type_ids: [],
        employee_ids: [],
        overtime_threshold: 45,
        overtime_multiplier: 1.5,
    });
    const [resultCount, setResultCount] = useState(0);

    useEffect(() => {
        Promise.all([
            departmentApi.list(),
            shiftTypeApi.list(),
        ]).then(([deptRes, stRes]) => {
            setDepartments(deptRes.data);
            setShiftTypes(stRes.data);
            setSelectedShiftTypes(stRes.data.map((st: ShiftType) => st.ID));
        }).catch(console.error);
    }, []);

    // Fetch employees when department changes
    useEffect(() => {
        if (params.department_id > 0) {
            setLoadingEmployees(true);
            employeeApi.list().then(res => {
                const deptEmployees = (res.data || []).filter(
                    (e: Employee) =>
                        e.DepartmentID === params.department_id &&
                        e.IsActive &&
                        (e.IsShiftWorker !== false) // Include unless explicitly false
                );
                setEmployees(deptEmployees);
                // Select all by default
                setSelectedEmployees(deptEmployees.map((e: Employee) => e.ID));
            }).catch(console.error).finally(() => setLoadingEmployees(false));
        } else {
            setEmployees([]);
            setSelectedEmployees([]);
        }
    }, [params.department_id]);

    const MONTHS_TR = [
        'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
        'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
    ];

    const toggleShiftType = (id: number) => {
        setSelectedShiftTypes(prev =>
            prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]
        );
    };

    const toggleEmployee = (id: number) => {
        setSelectedEmployees(prev =>
            prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]
        );
    };

    const toggleAllEmployees = () => {
        if (selectedEmployees.length === employees.length) {
            setSelectedEmployees([]);
        } else {
            setSelectedEmployees(employees.map(e => e.ID));
        }
    };

    const handleGenerate = async () => {
        if (!params.department_id) {
            alert('Lütfen bir bölüm seçin.');
            return;
        }
        if (selectedShiftTypes.length === 0) {
            alert('Lütfen en az bir nöbet tipi seçin.');
            return;
        }
        if (selectedEmployees.length === 0) {
            alert('Lütfen en az bir personel seçin.');
            return;
        }
        setLoading(true);
        try {
            const res = await scheduleApi.generate({
                ...params,
                shift_type_ids: selectedShiftTypes,
                employee_ids: selectedEmployees,
            });
            setResultCount(res.data?.length || 0);
            setStep(3);
        } catch (err) {
            alert('Oluşturma hatası: ' + err);
        } finally {
            setLoading(false);
        }
    };

    const selectedDept = departments.find(d => d.ID === params.department_id);

    return (
        <div className="max-w-2xl mx-auto py-8">
            {/* Stepper */}
            <div className="mb-12 flex justify-between items-center relative">
                <div className="absolute top-1/2 left-0 w-full h-0.5 bg-white/5 -z-10"></div>
                {[
                    { num: 1, label: 'Başlangıç' },
                    { num: 2, label: 'Parametreler' },
                    { num: 3, label: 'Sonuç' },
                ].map(({ num, label }) => (
                    <div key={num} className="flex flex-col items-center gap-2">
                        <div
                            className={`w-10 h-10 rounded-full flex items-center justify-center font-bold border-2 transition-all duration-500 ${step >= num
                                ? 'bg-gradient-to-br from-blue-500 to-blue-600 border-blue-500 text-white shadow-lg shadow-blue-500/20'
                                : 'bg-[var(--bg-card)] border-white/10 text-gray-500'
                                }`}
                        >
                            {step > num ? <CheckCircle2 className="w-5 h-5" /> : num}
                        </div>
                        <span className={`text-xs font-medium ${step >= num ? 'text-blue-400' : 'text-gray-600'}`}>
                            {label}
                        </span>
                    </div>
                ))}
            </div>

            {step === 1 && (
                <div className="space-y-6 text-center animate-scale-in">
                    <div className="w-20 h-20 mx-auto rounded-2xl bg-gradient-to-br from-blue-500/10 to-purple-500/10 flex items-center justify-center border border-blue-500/20">
                        <Wand2 className="w-10 h-10 text-blue-400 opacity-60" />
                    </div>
                    <h3 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                        Akıllı Nöbet Planlayıcı
                    </h3>
                    <p className="text-gray-400 max-w-md mx-auto">
                        Algoritmamız seçtiğiniz personeli adil şekilde dağıtarak,
                        ardışık nöbet olmadan en ideal listeyi oluşturur.
                    </p>
                    <button
                        onClick={() => setStep(2)}
                        className="btn-primary text-lg px-10 py-3 mx-auto transition-all hover:scale-105"
                    >
                        Başlayalım <ChevronRight className="w-5 h-5" />
                    </button>
                </div>
            )}

            {step === 2 && (
                <div className="glass-card p-8 animate-slide-up">
                    <h3 className="text-xl font-bold mb-6">Planlama Parametreleri</h3>
                    <div className="grid grid-cols-2 gap-6">
                        <div className="space-y-2 col-span-2">
                            <label className="text-sm text-gray-400 font-medium">Bölüm *</label>
                            <select
                                className="glass-input w-full"
                                value={params.department_id}
                                onChange={(e) => setParams({ ...params, department_id: Number(e.target.value) })}
                                required
                            >
                                <option value={0} disabled>Bölüm Seçin</option>
                                {departments.map((d) => (
                                    <option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>
                                ))}
                            </select>
                        </div>

                        {/* Shift Type Selection */}
                        <div className="space-y-3 col-span-2">
                            <label className="flex items-center justify-between font-semibold">
                                Çalışma Tipleri * <span className="text-xs text-gray-600">({selectedShiftTypes.length}/{shiftTypes.length} seçili)</span>
                            </label>
                            {shiftTypes.length === 0 ? (
                                <div className="text-sm text-red-400 bg-red-500/10 p-3 rounded-lg border border-red-500/20 flex items-center gap-2">
                                    <AlertCircle className="w-4 h-4" />
                                    Henüz çalışma tipi tanımlı değil. Önce sol menüden "Çalışma Tipleri" bölümüne gidip en az bir tip ekleyin.
                                </div>
                            ) : (
                                <div className="flex flex-wrap gap-2">
                                    {shiftTypes.map((st) => {
                                        const isSelected = selectedShiftTypes.includes(st.ID);
                                        return (
                                            <button
                                                key={st.ID}
                                                type="button"
                                                onClick={() => toggleShiftType(st.ID)}
                                                className="flex items-center gap-2 px-3 py-2 rounded-lg border transition-all text-sm"
                                                style={{
                                                    backgroundColor: isSelected ? (st.Color || '#3b82f6') + '20' : 'transparent',
                                                    borderColor: isSelected ? (st.Color || '#3b82f6') + '50' : 'rgba(255,255,255,0.08)',
                                                    color: isSelected ? (st.Color || '#3b82f6') : '#9ca3af',
                                                }}
                                            >
                                                <div
                                                    className="w-3 h-3 rounded-sm border-2 flex items-center justify-center"
                                                    style={{
                                                        borderColor: isSelected ? st.Color || '#3b82f6' : '#6b7280',
                                                        backgroundColor: isSelected ? st.Color || '#3b82f6' : 'transparent',
                                                    }}
                                                >
                                                    {isSelected && <span className="text-white text-[8px] font-bold">✓</span>}
                                                </div>
                                                {st.Name}
                                                <span className="text-[10px] opacity-60">{st.StartTime}–{st.EndTime}</span>
                                            </button>
                                        );
                                    })}
                                </div>
                            )}
                        </div>

                        {/* Employee Selection */}
                        <div className="space-y-3 col-span-2">
                            <div className="flex items-center justify-between">
                                <label className="text-sm text-gray-400 font-medium flex items-center gap-2">
                                    <Users className="w-4 h-4" />
                                    Personel * <span className="text-xs text-gray-600">({selectedEmployees.length}/{employees.length} seçili)</span>
                                </label>
                                {employees.length > 0 && (
                                    <button
                                        type="button"
                                        onClick={toggleAllEmployees}
                                        className="text-xs text-blue-400 hover:text-blue-300 transition-colors"
                                    >
                                        {selectedEmployees.length === employees.length ? 'Hiçbirini Seçme' : 'Tümünü Seç'}
                                    </button>
                                )}
                            </div>
                            {params.department_id === 0 ? (
                                <div className="text-sm text-gray-500 bg-white/[0.02] border border-white/5 rounded-lg px-4 py-3">
                                    Önce bir bölüm seçin.
                                </div>
                            ) : loadingEmployees ? (
                                <div className="flex items-center gap-2 text-sm text-gray-400 py-2">
                                    <Loader2 className="w-4 h-4 animate-spin" /> Personel yükleniyor...
                                </div>
                            ) : employees.length === 0 ? (
                                <div className="text-sm text-amber-400/80 bg-amber-500/5 border border-amber-500/10 rounded-lg px-4 py-3">
                                    Bu bölümde aktif personel bulunamadı.
                                </div>
                            ) : (
                                <div className="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto pr-1">
                                    {employees.map((emp) => {
                                        const isSelected = selectedEmployees.includes(emp.ID);
                                        return (
                                            <button
                                                key={emp.ID}
                                                type="button"
                                                onClick={() => toggleEmployee(emp.ID)}
                                                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg border text-sm text-left transition-all ${isSelected
                                                    ? 'bg-blue-500/10 border-blue-500/30 text-white'
                                                    : 'bg-transparent border-white/[0.06] text-gray-500 hover:border-white/10'
                                                    }`}
                                            >
                                                <div
                                                    className="w-4 h-4 rounded border-2 flex items-center justify-center flex-shrink-0"
                                                    style={{
                                                        borderColor: isSelected ? '#3b82f6' : '#6b7280',
                                                        backgroundColor: isSelected ? '#3b82f6' : 'transparent',
                                                    }}
                                                >
                                                    {isSelected && <span className="text-white text-[8px] font-bold">✓</span>}
                                                </div>
                                                <div className="truncate">
                                                    <div className="font-medium">{emp.FirstName} {emp.LastName}</div>
                                                    <div className="text-[10px] text-gray-500">{emp.Title?.Name}</div>
                                                </div>
                                            </button>
                                        );
                                    })}
                                </div>
                            )}
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm text-gray-400 font-medium">Yıl</label>
                            <input
                                type="number"
                                className="glass-input w-full"
                                value={params.year}
                                onChange={(e) => setParams({ ...params, year: Number(e.target.value) })}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400 font-medium">Ay</label>
                            <select
                                className="glass-input w-full"
                                value={params.month}
                                onChange={(e) => setParams({ ...params, month: Number(e.target.value) })}
                            >
                                {MONTHS_TR.map((name, i) => (
                                    <option key={i + 1} value={i + 1}>{name}</option>
                                ))}
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400 font-medium">Ek Mesai Eşiği (Saat)</label>
                            <input
                                type="number"
                                className="glass-input w-full"
                                value={params.overtime_threshold}
                                onChange={(e) => setParams({ ...params, overtime_threshold: Number(e.target.value) })}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400 font-medium">Mesai Çarpanı (x)</label>
                            <input
                                type="number"
                                step="0.1"
                                className="glass-input w-full"
                                value={params.overtime_multiplier}
                                onChange={(e) => setParams({ ...params, overtime_multiplier: Number(e.target.value) })}
                            />
                        </div>
                    </div>
                    <div className="flex gap-3 mt-8">
                        <button onClick={() => setStep(1)} className="btn-ghost">Geri</button>
                        <button
                            onClick={handleGenerate}
                            disabled={loading || !params.department_id || selectedShiftTypes.length === 0 || selectedEmployees.length === 0}
                            className="btn-success flex-1 justify-center py-3 text-lg"
                        >
                            {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : <Wand2 className="w-5 h-5" />}
                            {loading ? 'Optimize ediliyor...' : 'Sihirbazı Çalıştır'}
                        </button>
                    </div>
                </div>
            )}

            {step === 3 && (
                <div className="space-y-6 text-center animate-scale-in">
                    <div className="w-20 h-20 bg-green-900/20 rounded-full flex items-center justify-center mx-auto border-2 border-green-500/50 shadow-lg shadow-green-500/10">
                        <CheckCircle2 className="w-10 h-10 text-green-400" />
                    </div>
                    <h3 className="text-2xl font-bold text-white">Çizelge Başarıyla Oluşturuldu!</h3>
                    <p className="text-gray-400">
                        {selectedDept && <span className="text-blue-400 font-semibold">{selectedDept.Floor}. Kat - {selectedDept.Name}</span>}
                        {' '}bölümünden{' '}
                        <span className="text-purple-400 font-semibold">{selectedEmployees.length} personel</span> için{' '}
                        <span className="text-green-400 font-bold text-lg">{resultCount}</span> adet nöbet ataması optimize edilerek kaydedildi.
                    </p>
                    <div className="flex gap-4 justify-center mt-8">
                        <button onClick={() => setStep(1)} className="btn-ghost">
                            Yeni Plan
                        </button>
                        <button
                            onClick={() => onNavigate?.('schedule')}
                            className="btn-primary"
                        >
                            Takvime Git
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ScheduleWizard;
