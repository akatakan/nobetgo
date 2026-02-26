import React, { useState, useEffect } from 'react';
import { ShieldCheck, Loader2, CheckCircle2, XCircle, Clock, CalendarOff, RefreshCw } from 'lucide-react';
import { approvalApi } from '../services/api';
import type { TimeEntry, Leave, AuditLog } from '../types';

const ApprovalManager: React.FC = () => {
    const [pendingEntries, setPendingEntries] = useState<TimeEntry[]>([]);
    const [pendingLeaves, setPendingLeaves] = useState<Leave[]>([]);
    const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [activeTab, setActiveTab] = useState<'pending' | 'logs'>('pending');

    const fetchPending = async () => {
        setLoading(true);
        try {
            const res = await approvalApi.getPending();
            setPendingEntries(res.data?.time_entries || []);
            setPendingLeaves(res.data?.leaves || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const fetchLogs = async () => {
        try {
            const res = await approvalApi.getAuditLogs('time_entry');
            setAuditLogs(res.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchPending();
        fetchLogs();
    }, []);

    const handleApproveEntry = async (id: number) => {
        setActionLoading(`approve-entry-${id}`);
        try {
            await approvalApi.approveTimeEntry(id, 1);
            fetchPending();
            fetchLogs();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Onay hatası');
        } finally {
            setActionLoading(null);
        }
    };

    const handleRejectEntry = async (id: number) => {
        setActionLoading(`reject-entry-${id}`);
        try {
            await approvalApi.rejectTimeEntry(id, 1);
            fetchPending();
            fetchLogs();
        } catch (err: any) {
            alert(err?.response?.data?.error || 'Red hatası');
        } finally {
            setActionLoading(null);
        }
    };

    const formatTime = (d: string) => new Date(d).toLocaleString('tr-TR', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' });

    const totalPending = pendingEntries.length + pendingLeaves.length;

    return (
        <div className="space-y-6 animate-fade-in">
            {/* Header */}
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold flex items-center gap-2">
                    <ShieldCheck className="w-5 h-5 text-blue-400" /> Onay Merkezi
                </h2>
                <div className="flex items-center gap-2">
                    <button onClick={() => { fetchPending(); fetchLogs(); }} className="btn-ghost text-xs py-2 px-3 flex items-center gap-2">
                        <RefreshCw className="w-3.5 h-3.5" /> Yenile
                    </button>
                </div>
            </div>

            {/* Tab Switch */}
            <div className="flex gap-2">
                <button
                    onClick={() => setActiveTab('pending')}
                    className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${activeTab === 'pending' ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20' : 'text-gray-400 hover:bg-white/5'}`}
                >
                    Bekleyenler
                    {totalPending > 0 && (
                        <span className="ml-2 bg-amber-500/20 text-amber-400 text-xs px-1.5 py-0.5 rounded-full font-bold">{totalPending}</span>
                    )}
                </button>
                <button
                    onClick={() => setActiveTab('logs')}
                    className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${activeTab === 'logs' ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20' : 'text-gray-400 hover:bg-white/5'}`}
                >
                    Denetim İzi
                </button>
            </div>

            {activeTab === 'pending' && (
                <>
                    {loading ? (
                        <div className="flex items-center justify-center py-20 text-gray-500">
                            <Loader2 className="w-6 h-6 animate-spin mr-3" /> Yükleniyor...
                        </div>
                    ) : totalPending === 0 ? (
                        <div className="glass-card p-12 text-center text-gray-500">
                            <CheckCircle2 className="w-12 h-12 mx-auto mb-3 opacity-30" />
                            <p>Bekleyen onay yok. Tüm kayıtlar güncel.</p>
                        </div>
                    ) : (
                        <div className="space-y-5">
                            {/* Pending Time Entries */}
                            {pendingEntries.length > 0 && (
                                <div className="glass-card overflow-hidden">
                                    <div className="px-5 py-3 bg-white/[0.02] border-b border-white/5 flex items-center gap-2">
                                        <Clock className="w-4 h-4 text-blue-400" />
                                        <span className="text-sm font-semibold text-gray-300">Puantaj Kayıtları</span>
                                        <span className="text-xs text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-full ml-2">{pendingEntries.length}</span>
                                    </div>
                                    <div className="divide-y divide-white/[0.03]">
                                        {pendingEntries.map(entry => (
                                            <div key={entry.ID} className="px-5 py-3 flex items-center gap-4 animate-slide-up">
                                                <div className="w-8 h-8 rounded-lg bg-amber-500/10 flex items-center justify-center">
                                                    <Clock className="w-4 h-4 text-amber-400" />
                                                </div>
                                                <div className="flex-1">
                                                    <div className="font-medium text-sm">
                                                        {entry.employee?.FirstName} {entry.employee?.LastName}
                                                    </div>
                                                    <div className="text-xs text-gray-500">
                                                        {formatTime(entry.clock_in)}
                                                        {entry.clock_out && ` → ${formatTime(entry.clock_out)}`}
                                                        <span className="ml-2 text-gray-600">({entry.source})</span>
                                                    </div>
                                                </div>
                                                <div className="flex gap-2">
                                                    <button
                                                        onClick={() => handleApproveEntry(entry.ID)}
                                                        disabled={actionLoading === `approve-entry-${entry.ID}`}
                                                        className="btn-ghost px-3 py-1.5 text-xs text-emerald-400 hover:bg-emerald-500/10 flex items-center gap-1"
                                                    >
                                                        {actionLoading === `approve-entry-${entry.ID}` ? <Loader2 className="w-3 h-3 animate-spin" /> : <CheckCircle2 className="w-3.5 h-3.5" />}
                                                        Onayla
                                                    </button>
                                                    <button
                                                        onClick={() => handleRejectEntry(entry.ID)}
                                                        disabled={actionLoading === `reject-entry-${entry.ID}`}
                                                        className="btn-ghost px-3 py-1.5 text-xs text-red-400 hover:bg-red-500/10 flex items-center gap-1"
                                                    >
                                                        {actionLoading === `reject-entry-${entry.ID}` ? <Loader2 className="w-3 h-3 animate-spin" /> : <XCircle className="w-3.5 h-3.5" />}
                                                        Reddet
                                                    </button>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Pending Leaves */}
                            {pendingLeaves.length > 0 && (
                                <div className="glass-card overflow-hidden">
                                    <div className="px-5 py-3 bg-white/[0.02] border-b border-white/5 flex items-center gap-2">
                                        <CalendarOff className="w-4 h-4 text-purple-400" />
                                        <span className="text-sm font-semibold text-gray-300">İzin Talepleri</span>
                                        <span className="text-xs text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-full ml-2">{pendingLeaves.length}</span>
                                    </div>
                                    <div className="divide-y divide-white/[0.03]">
                                        {pendingLeaves.map(leave => (
                                            <div key={leave.ID} className="px-5 py-3 flex items-center gap-4 animate-slide-up">
                                                <div className="w-8 h-8 rounded-lg bg-purple-500/10 flex items-center justify-center">
                                                    <CalendarOff className="w-4 h-4 text-purple-400" />
                                                </div>
                                                <div className="flex-1">
                                                    <div className="font-medium text-sm">
                                                        {leave.employee?.FirstName} {leave.employee?.LastName}
                                                    </div>
                                                    <div className="text-xs text-gray-500">
                                                        {leave.leave_type?.name} · {leave.start_date?.split('T')[0]} → {leave.end_date?.split('T')[0]} · {leave.total_days} gün
                                                    </div>
                                                    {leave.reason && <div className="text-xs text-gray-600 mt-0.5 italic">{leave.reason}</div>}
                                                </div>
                                                <span className="text-xs text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-md">Bekliyor</span>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </>
            )}

            {activeTab === 'logs' && (
                <div className="glass-card overflow-hidden">
                    <div className="px-5 py-3 bg-white/[0.02] border-b border-white/5">
                        <span className="text-sm font-semibold text-gray-300">Son İşlemler</span>
                    </div>
                    {auditLogs.length === 0 ? (
                        <div className="p-8 text-center text-gray-500 text-sm">Denetim kaydı bulunamadı.</div>
                    ) : (
                        <div className="divide-y divide-white/[0.03] max-h-[500px] overflow-y-auto">
                            {auditLogs.map(log => (
                                <div key={log.ID} className="px-5 py-3 flex items-center gap-4 text-sm animate-slide-up">
                                    <div className={`w-2 h-2 rounded-full flex-shrink-0 ${log.action === 'approve' ? 'bg-emerald-400' : log.action === 'reject' ? 'bg-red-400' : 'bg-blue-400'}`} />
                                    <div className="flex-1">
                                        <span className="text-gray-300">{log.entity_type}</span>
                                        <span className="text-gray-600 mx-1">#{log.entity_id}</span>
                                        <span className={`font-medium ${log.action === 'approve' ? 'text-emerald-400' : log.action === 'reject' ? 'text-red-400' : 'text-blue-400'}`}>
                                            {log.action}
                                        </span>
                                        {log.field_name && (
                                            <span className="text-gray-500 ml-1">
                                                ({log.field_name}: {log.old_value} → {log.new_value})
                                            </span>
                                        )}
                                    </div>
                                    <div className="text-xs text-gray-600">{new Date(log.CreatedAt).toLocaleString('tr-TR')}</div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default ApprovalManager;
