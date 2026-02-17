import React, { useState, useEffect } from 'react';
import { Plus, Clock, Palette, Pencil, Trash2, Loader2, AlertTriangle } from 'lucide-react';
import { shiftTypeApi } from '../services/api';
import type { ShiftType, ShiftTypeFormData } from '../types';

const emptyForm: ShiftTypeFormData = {
    Name: '',
    StartTime: '08:00',
    EndTime: '08:00',
    Color: '#3b82f6',
    Description: '',
};

const ShiftTypeManager: React.FC = () => {
    const [shifts, setShifts] = useState<ShiftType[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [formData, setFormData] = useState<ShiftTypeFormData>({ ...emptyForm });
    const [confirmDelete, setConfirmDelete] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);

    const fetchShifts = async () => {
        setLoading(true);
        try {
            const res = await shiftTypeApi.list();
            setShifts(res.data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchShifts();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingId) {
                await shiftTypeApi.update(editingId, formData);
            } else {
                await shiftTypeApi.create(formData);
            }
            setShowForm(false);
            setEditingId(null);
            setFormData({ ...emptyForm });
            fetchShifts();
        } catch (err) {
            alert('Hata oluştu');
        } finally {
            setSaving(false);
        }
    };

    const handleEdit = (s: ShiftType) => {
        setEditingId(s.ID);
        setFormData({
            Name: s.Name,
            StartTime: s.StartTime,
            EndTime: s.EndTime,
            Color: s.Color,
            Description: s.Description,
        });
        setShowForm(true);
    };

    const handleDelete = async (id: number) => {
        try {
            await shiftTypeApi.delete(id);
            setConfirmDelete(null);
            fetchShifts();
        } catch (err) {
            alert('Silme hatası');
        }
    };

    const cancelForm = () => {
        setShowForm(false);
        setEditingId(null);
        setFormData({ ...emptyForm });
    };

    // Calculate duration string
    const calcDuration = (start: string, end: string): string => {
        const [sh, sm] = start.split(':').map(Number);
        const [eh, em] = end.split(':').map(Number);
        let diff = (eh * 60 + em) - (sh * 60 + sm);
        if (diff <= 0) diff += 24 * 60; // overnight
        const hours = Math.floor(diff / 60);
        const mins = diff % 60;
        return mins > 0 ? `${hours}s ${mins}dk` : `${hours} saat`;
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">
                    Çalışma Tipleri
                    <span className="text-sm font-normal text-gray-500 ml-2">({shifts.length})</span>
                </h2>
                <button
                    onClick={() => { cancelForm(); setShowForm(!showForm); }}
                    className="btn-primary"
                >
                    <Plus className="w-4 h-4" />
                    Yeni Nöbet Tipi
                </button>
            </div>

            {/* Form */}
            {showForm && (
                <form onSubmit={handleSubmit} className="glass-card p-6 animate-slide-down">
                    <h3 className="text-lg font-semibold mb-4">
                        {editingId ? 'Nöbet Tipi Düzenle' : 'Yeni Nöbet Tipi'}
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Nöbet Adı *</label>
                            <input
                                placeholder="Örn: 24 Saatlik"
                                className="glass-input w-full"
                                value={formData.Name}
                                onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Renk</label>
                            <div className="flex items-center gap-2 glass-input">
                                <Palette className="w-4 h-4 text-gray-400" />
                                <input
                                    type="color"
                                    className="bg-transparent border-none w-8 h-6 p-0 cursor-pointer"
                                    value={formData.Color}
                                    onChange={(e) => setFormData({ ...formData, Color: e.target.value })}
                                />
                                <span className="text-xs text-gray-500">{formData.Color}</span>
                            </div>
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Başlangıç *</label>
                            <input
                                type="time"
                                className="glass-input w-full"
                                value={formData.StartTime}
                                onChange={(e) => setFormData({ ...formData, StartTime: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1">
                            <label className="text-xs text-gray-400 font-medium">Bitiş *</label>
                            <input
                                type="time"
                                className="glass-input w-full"
                                value={formData.EndTime}
                                onChange={(e) => setFormData({ ...formData, EndTime: e.target.value })}
                                required
                            />
                        </div>
                        <div className="space-y-1 md:col-span-2">
                            <label className="text-xs text-gray-400 font-medium">Açıklama</label>
                            <textarea
                                placeholder="Açıklama (opsiyonel)"
                                className="glass-input w-full resize-none"
                                rows={2}
                                value={formData.Description}
                                onChange={(e) => setFormData({ ...formData, Description: e.target.value })}
                            />
                        </div>
                    </div>
                    <div className="flex justify-end gap-3 mt-6">
                        <button type="button" onClick={cancelForm} className="btn-ghost">İptal</button>
                        <button type="submit" disabled={saving} className="btn-success">
                            {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                            {editingId ? 'Güncelle' : 'Kaydet'}
                        </button>
                    </div>
                </form>
            )}

            {/* Shift cards */}
            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="skeleton h-28 rounded-xl" />
                    ))}
                </div>
            ) : shifts.length === 0 ? (
                <div className="text-center py-16 text-gray-500">
                    <Clock className="w-12 h-12 mx-auto mb-3 opacity-30" />
                    <p>Henüz nöbet tipi eklenmemiş.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {shifts.map((s, idx) => (
                        <div
                            key={s.ID}
                            className="glass-card p-5 group hover:border-blue-500/30 transition-all duration-300 animate-slide-up"
                            style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
                        >
                            <div className="flex items-start gap-4">
                                <div
                                    className="w-11 h-11 rounded-xl flex items-center justify-center flex-shrink-0"
                                    style={{ backgroundColor: s.Color + '20' }}
                                >
                                    <Clock className="w-5 h-5" style={{ color: s.Color }} />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="font-semibold text-base flex items-center gap-2">
                                        <span className="truncate">{s.Name}</span>
                                        <span
                                            className="w-2.5 h-2.5 rounded-full flex-shrink-0"
                                            style={{ backgroundColor: s.Color }}
                                        />
                                    </div>
                                    <div className="text-gray-400 text-sm mt-0.5">
                                        {s.StartTime} — {s.EndTime}
                                    </div>
                                    <div className="text-xs text-gray-500 mt-1">
                                        {calcDuration(s.StartTime, s.EndTime)}
                                    </div>
                                    {s.Description && (
                                        <div className="text-xs text-gray-500 mt-1 truncate">{s.Description}</div>
                                    )}
                                </div>
                                <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                                    <button
                                        onClick={() => handleEdit(s)}
                                        className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-blue-400 transition-colors"
                                    >
                                        <Pencil className="w-3.5 h-3.5" />
                                    </button>
                                    <button
                                        onClick={() => setConfirmDelete(s.ID)}
                                        className="p-1.5 rounded-lg hover:bg-white/5 text-gray-400 hover:text-red-400 transition-colors"
                                    >
                                        <Trash2 className="w-3.5 h-3.5" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Delete Dialog */}
            {confirmDelete !== null && (
                <div className="dialog-overlay" onClick={() => setConfirmDelete(null)}>
                    <div className="dialog-content" onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-3 mb-4">
                            <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                                <AlertTriangle className="w-5 h-5 text-red-400" />
                            </div>
                            <div>
                                <h4 className="font-semibold">Nöbet Tipini Sil</h4>
                                <p className="text-sm text-gray-400">Bu işlem geri alınamaz.</p>
                            </div>
                        </div>
                        <div className="flex justify-end gap-3">
                            <button onClick={() => setConfirmDelete(null)} className="btn-ghost">İptal</button>
                            <button onClick={() => handleDelete(confirmDelete)} className="btn-danger">Sil</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ShiftTypeManager;
