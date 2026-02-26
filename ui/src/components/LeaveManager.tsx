import React, { useState, useEffect } from 'react';
import { CalendarOff, Plus, ChevronLeft, ChevronRight, Loader2, CheckCircle2, XCircle, Tag } from 'lucide-react';
import { leaveApi, leaveTypeApi, employeeApi, departmentApi } from '../services/api';
import type { Leave, LeaveType, Employee, Department, LeaveBalance } from '../types';

const MONTHS_TR = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
];

const STATUS_MAP: Record<string, { label: string; color: string; icon: React.ReactNode }> = {
    pending: { label: 'Bekliyor', color: 'text-amber-400 bg-amber-500/10', icon: <Loader2 className="w-3 h-3" /> },
    approved: { label: 'Onaylandı', color: 'text-emerald-400 bg-emerald-500/10', icon: <CheckCircle2 className="w-3 h-3" /> },
    rejected: { label: 'Reddedildi', color: 'text-red-400 bg-red-500/10', icon: <XCircle className="w-3 h-3" /> },
};

const LeaveManager: React.FC = () => {
    const now = new Date();
    const [month, setMonth] = useState(now.getMonth() + 1);
    const [year, setYear] = useState(now.getFullYear());
    const [leaves, setLeaves] = useState<Leave[]>([]);
    const [leaveTypes, setLeaveTypes] = useState<LeaveType[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [selectedDept, setSelectedDept] = useState<number>(0);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    // Request form
    const [showRequestForm, setShowRequestForm] = useState(false);
    const [requestForm, setRequestForm] = useState({
        employee_id: 0,
        leave_type_id: 0,
        start_date: '',
        end_date: '',
        reason: '',
    });

    // Leave type form
    const [showTypeForm, setShowTypeForm] = useState(false);
    const [typeForm, setTypeForm] = useState({ name: '', default_days: 14, is_paid: true, requires_approval: true, color: '#3B82F6' });

    // Balance view
    const [balanceEmployee, setBalanceEmployee] = useState<number>(0);
    const [balances, setBalances] = useState<LeaveBalance[]>([]);

    useEffect(() => {
        Promise.all([
            leaveTypeApi.list(),
            employeeApi.list(),
            departmentApi.list(),
        ]).then(([ltRes, empRes, deptRes]) => {
            setLeaveTypes(ltRes.data || []);
            setEmployees(empRes.data || []);
            setDepartments(deptRes.data || []);
            if (deptRes.data.length > 0 && selectedDept === 0) setSelectedDept(deptRes.data[0].ID);
        }).catch(console.error);
    }, []);

    const fetchLeaves = async () => {
        setLoading(true);
        try {
            const startDate = `${year}-${String(month).padStart(2, '0')}-01`;
            const endMonth = month === 12 ? 1 : month + 1;
            const endYear = month === 12 ? year + 1 : year;
            const endDate = `${endYear}-${String(endMonth).padStart(2, '0')}-01`;

            const params: Record<string, string | number> = { start: startDate, end: endDate };
            if (selectedDept) params.department_id = selectedDept;

            const res = await leaveApi.list(params);
            setLeaves(res.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (selectedDept) fetchLeaves();
    }, [month, year, selectedDept]);

    // Fetch balance
    useEffect(() => {
        if (balanceEmployee > 0) {
            leaveApi.getBalance(balanceEmployee, year).then(res => {
                setBalances(res.data || []);
            }).catch(console.error);
        }
    }, [balanceEmployee, year]);

    const prevMonth = () => { if (month === 1) { setMonth(12); setYear(year - 1); } else setMonth(month - 1); };
    const nextMonth = () => { if (month === 12) { setMonth(1); setYear(year + 1); } else setMonth(month + 1); };

    const handleRequestLeave = async () => {
        if (!requestForm.employee_id || !requestForm.leave_type_id || !requestForm.start_date || !requestForm.end_date) {
            alert('Tüm alanları doldurun.'); return;
        }
        setActionLoading('request');
        try {
            await leaveApi.request(requestForm);
            setShowRequestForm(false);
            setRequestForm({ employee_id: 0, leave_type_id: 0, start_date: '', end_date: '', reason: '' });
            fetchLeaves();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'İzin talebi hatası');
        } finally {
            setActionLoading(null);
        }
    };

    const handleApprove = async (id: number) => {
        setActionLoading(`approve-${id}`);
        try {
            await leaveApi.approve(id, 1); // TODO: real approver from auth
            fetchLeaves();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Onay hatası');
        } finally {
            setActionLoading(null);
        }
    };

    const handleReject = async (id: number) => {
        setActionLoading(`reject-${id}`);
        try {
            await leaveApi.reject(id, 1);
            fetchLeaves();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Red hatası');
        } finally {
            setActionLoading(null);
        }
    };

    const handleCreateType = async () => {
        if (!typeForm.name) { alert('İzin türü adı gerekli.'); return; }
        try {
            await leaveTypeApi.create(typeForm);
            setShowTypeForm(false);
            setTypeForm({ name: '', default_days: 14, is_paid: true, requires_approval: true, color: '#3B82F6' });
            const res = await leaveTypeApi.list();
            setLeaveTypes(res.data || []);
        } catch (err) { console.error(err); }
    };

    const deptEmployees = employees.filter(e => !selectedDept || e.DepartmentID === selectedDept);

    // Stats
    const pendingCount = leaves.filter(l => l.status === 'pending').length;
    const approvedCount = leaves.filter(l => l.status === 'approved').length;
    const totalDays = leaves.filter(l => l.status === 'approved').reduce((s, l) => s + l.total_days, 0);

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center flex-wrap gap-3">
                <div className="flex items-center gap-4">
                    <h2 className="text-xl font-semibold">İzin Yönetimi</h2>
                    <select className="glass-input py-1.5 px-3 text-sm" value={selectedDept} onChange={(e) => setSelectedDept(Number(e.target.value))}>
                        <option value={0} disabled>Bölüm Seçin</option>
                        {departments.map((d) => (<option key={d.ID} value={d.ID}>{d.Floor}. Kat - {d.Name}</option>))}
                    </select>
                </div>
                <div className="flex items-center gap-2">
                    <button onClick={() => setShowTypeForm(true)} className="btn-ghost text-xs py-2 px-3 flex items-center gap-2">
                        <Tag className="w-3.5 h-3.5" /> İzin Türü
                    </button>
                    <button onClick={() => setShowRequestForm(true)} className="btn-primary text-xs py-2 px-3 flex items-center gap-2">
                        <Plus className="w-3.5 h-3.5" /> İzin Talebi
                    </button>
                    <div className="h-6 w-px bg-white/10 mx-1"></div>
                    <button onClick={prevMonth} className="btn-ghost p-2"><ChevronLeft className="w-4 h-4" /></button>
                    <span className="min-w-[160px] text-center font-semibold text-lg">{MONTHS_TR[month - 1]} {year}</span>
                    <button onClick={nextMonth} className="btn-ghost p-2"><ChevronRight className="w-4 h-4" /></button>
                </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="glass-card p-4 bg-gradient-to-br from-amber-500/5 to-amber-500/10 border-amber-500/10">
                    <div className="text-xs text-gray-400 mb-1">Bekleyen Talepler</div>
                    <div className="text-2xl font-bold text-amber-400">{pendingCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-emerald-500/5 to-emerald-500/10 border-emerald-500/10">
                    <div className="text-xs text-gray-400 mb-1">Onaylanan</div>
                    <div className="text-2xl font-bold text-emerald-400">{approvedCount}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-blue-500/5 to-blue-500/10 border-blue-500/10">
                    <div className="text-xs text-gray-400 mb-1">Toplam İzin Günü</div>
                    <div className="text-2xl font-bold text-blue-400">{totalDays}</div>
                </div>
                <div className="glass-card p-4 bg-gradient-to-br from-purple-500/5 to-purple-500/10 border-purple-500/10">
                    <div className="text-xs text-gray-400 mb-1">İzin Türleri</div>
                    <div className="text-2xl font-bold text-purple-400">{leaveTypes.length}</div>
                </div>
            </div>

            {/* Balance Quick View */}
            <div className="glass-card p-5 bg-gradient-to-br from-indigo-500/5 to-purple-500/5 border-indigo-500/10">
                <h3 className="text-sm font-semibold text-gray-300 mb-3">İzin Bakiyesi</h3>
                <div className="flex items-center gap-3">
                    <select className="glass-input py-1.5 px-3 text-sm flex-1" value={balanceEmployee} onChange={(e) => setBalanceEmployee(Number(e.target.value))}>
                        <option value={0}>Personel seçin...</option>
                        {deptEmployees.map(e => (<option key={e.ID} value={e.ID}>{e.FirstName} {e.LastName}</option>))}
                    </select>
                </div>
                {balances.length > 0 && (
                    <div className="mt-3 grid grid-cols-2 md:grid-cols-4 gap-3">
                        {balances.map(b => (
                            <div key={b.ID} className="bg-white/[0.03] rounded-xl p-3 border border-white/5">
                                <div className="text-xs text-gray-400">{b.leave_type?.name || 'İzin'}</div>
                                <div className="text-lg font-bold text-white mt-1">
                                    {b.remaining_days} <span className="text-xs text-gray-500 font-normal">/ {b.total_days} gün</span>
                                </div>
                                <div className="w-full h-1.5 bg-white/5 rounded-full mt-2 overflow-hidden">
                                    <div
                                        className="h-full bg-gradient-to-r from-blue-500 to-blue-400 rounded-full transition-all"
                                        style={{ width: `${Math.max(0, (b.remaining_days / b.total_days) * 100)}%` }}
                                    />
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Leave List */}
            {loading ? (
                <div className="flex items-center justify-center py-20 text-gray-500">
                    <Loader2 className="w-6 h-6 animate-spin mr-3" />Yükleniyor...
                </div>
            ) : leaves.length === 0 ? (
                <div className="glass-card p-12 text-center text-gray-500">
                    <CalendarOff className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Bu ay için izin kaydı bulunamadı.</p>
                </div>
            ) : (
                <div className="glass-card overflow-hidden">
                    <div className="grid grid-cols-7 gap-2 px-5 py-3 bg-white/[0.02] border-b border-white/5 text-xs font-semibold text-gray-400 uppercase">
                        <div className="col-span-2">Personel</div>
                        <div>İzin Türü</div>
                        <div>Başlangıç</div>
                        <div>Bitiş</div>
                        <div className="text-center">Gün</div>
                        <div className="text-center">Durum</div>
                    </div>
                    {leaves.map((leave, idx) => {
                        const statusInfo = STATUS_MAP[leave.status] || STATUS_MAP.pending;
                        return (
                            <div
                                key={leave.ID}
                                className="grid grid-cols-7 gap-2 px-5 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors items-center animate-slide-up"
                                style={{ animationDelay: `${idx * 30}ms`, animationFillMode: 'both' }}
                            >
                                <div className="col-span-2 flex items-center gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/15 to-purple-500/15 flex items-center justify-center text-xs font-bold text-blue-400">
                                        {leave.employee?.FirstName?.[0]}{leave.employee?.LastName?.[0]}
                                    </div>
                                    <div className="font-medium text-sm">{leave.employee?.FirstName} {leave.employee?.LastName}</div>
                                </div>
                                <div className="text-sm">
                                    <span
                                        className="px-2 py-0.5 rounded-md text-xs font-medium"
                                        style={{ backgroundColor: `${leave.leave_type?.color || '#3B82F6'}20`, color: leave.leave_type?.color || '#3B82F6' }}
                                    >
                                        {leave.leave_type?.name}
                                    </span>
                                </div>
                                <div className="text-sm text-gray-300">{leave.start_date?.split('T')[0]}</div>
                                <div className="text-sm text-gray-300">{leave.end_date?.split('T')[0]}</div>
                                <div className="text-center text-sm font-bold">{leave.total_days}</div>
                                <div className="text-center flex items-center justify-center gap-2">
                                    <span className={`text-[10px] px-2 py-0.5 rounded-md font-medium flex items-center gap-1 ${statusInfo.color}`}>
                                        {statusInfo.icon} {statusInfo.label}
                                    </span>
                                    {leave.status === 'pending' && (
                                        <div className="flex gap-1">
                                            <button
                                                onClick={() => handleApprove(leave.ID)}
                                                disabled={actionLoading === `approve-${leave.ID}`}
                                                className="btn-ghost p-1 text-emerald-400/60 hover:text-emerald-400"
                                                title="Onayla"
                                            >
                                                <CheckCircle2 className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleReject(leave.ID)}
                                                disabled={actionLoading === `reject-${leave.ID}`}
                                                className="btn-ghost p-1 text-red-400/60 hover:text-red-400"
                                                title="Reddet"
                                            >
                                                <XCircle className="w-4 h-4" />
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}

            {/* Request Form Modal */}
            {showRequestForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 animate-fade-in">
                    <div className="glass-card p-6 w-full max-w-md bg-[#1e293b]">
                        <h3 className="text-lg font-semibold mb-4">Yeni İzin Talebi</h3>
                        <div className="space-y-4">
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Personel</label>
                                <select className="glass-input w-full" value={requestForm.employee_id} onChange={(e) => setRequestForm({ ...requestForm, employee_id: Number(e.target.value) })}>
                                    <option value={0}>Seçin</option>
                                    {deptEmployees.map(e => (<option key={e.ID} value={e.ID}>{e.FirstName} {e.LastName}</option>))}
                                </select>
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">İzin Türü</label>
                                <select className="glass-input w-full" value={requestForm.leave_type_id} onChange={(e) => setRequestForm({ ...requestForm, leave_type_id: Number(e.target.value) })}>
                                    <option value={0}>Seçin</option>
                                    {leaveTypes.map(lt => (<option key={lt.ID} value={lt.ID}>{lt.name}</option>))}
                                </select>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Başlangıç</label>
                                    <input type="date" className="glass-input w-full" value={requestForm.start_date} onChange={(e) => setRequestForm({ ...requestForm, start_date: e.target.value })} />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Bitiş</label>
                                    <input type="date" className="glass-input w-full" value={requestForm.end_date} onChange={(e) => setRequestForm({ ...requestForm, end_date: e.target.value })} />
                                </div>
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">Açıklama</label>
                                <textarea className="glass-input w-full" rows={2} value={requestForm.reason} onChange={(e) => setRequestForm({ ...requestForm, reason: e.target.value })} placeholder="İsteğe bağlı" />
                            </div>
                        </div>
                        <div className="flex justify-end gap-3 mt-6">
                            <button onClick={() => setShowRequestForm(false)} className="btn-ghost">İptal</button>
                            <button onClick={handleRequestLeave} disabled={actionLoading === 'request'} className="btn-primary">
                                {actionLoading === 'request' ? <Loader2 className="w-4 h-4 animate-spin" /> : null} Talep Oluştur
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Leave Type Form Modal */}
            {showTypeForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 animate-fade-in">
                    <div className="glass-card p-6 w-full max-w-md bg-[#1e293b]">
                        <h3 className="text-lg font-semibold mb-4">Yeni İzin Türü</h3>
                        <div className="space-y-4">
                            <div className="space-y-1">
                                <label className="text-xs text-gray-400 font-medium">İzin Türü Adı</label>
                                <input className="glass-input w-full" value={typeForm.name} onChange={(e) => setTypeForm({ ...typeForm, name: e.target.value })} placeholder="Ör: Yıllık İzin" />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Varsayılan Gün</label>
                                    <input type="number" className="glass-input w-full" value={typeForm.default_days} onChange={(e) => setTypeForm({ ...typeForm, default_days: Number(e.target.value) })} />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-xs text-gray-400 font-medium">Renk</label>
                                    <input type="color" className="glass-input w-full h-10" value={typeForm.color} onChange={(e) => setTypeForm({ ...typeForm, color: e.target.value })} />
                                </div>
                            </div>
                            <div className="flex items-center gap-6">
                                <label className="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                                    <input type="checkbox" checked={typeForm.is_paid} onChange={(e) => setTypeForm({ ...typeForm, is_paid: e.target.checked })} className="w-4 h-4 rounded" /> Ücretli
                                </label>
                                <label className="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
                                    <input type="checkbox" checked={typeForm.requires_approval} onChange={(e) => setTypeForm({ ...typeForm, requires_approval: e.target.checked })} className="w-4 h-4 rounded" /> Onay Gerekli
                                </label>
                            </div>
                        </div>
                        <div className="flex justify-end gap-3 mt-6">
                            <button onClick={() => setShowTypeForm(false)} className="btn-ghost">İptal</button>
                            <button onClick={handleCreateType} className="btn-primary">Oluştur</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default LeaveManager;
