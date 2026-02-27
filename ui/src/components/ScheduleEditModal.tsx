import React, { useState, useEffect } from 'react';
import { X, Save, Trash2, Loader2, User, Clock } from 'lucide-react';
import { employeeApi, shiftTypeApi, scheduleApi } from '../services/api';
import type { Schedule, Employee, ShiftType } from '../types';

interface Props {
    isOpen: boolean;
    schedule: Schedule | null; // If null, we are creating a new one
    date: Date; // The date we are editing/creating for
    onClose: () => void;
    onSuccess: () => void;
}

const ScheduleEditModal: React.FC<Props> = ({ isOpen, schedule, date, onClose, onSuccess }) => {
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [shiftTypes, setShiftTypes] = useState<ShiftType[]>([]);

    // Form state
    const [formData, setFormData] = useState({
        employee_id: 0,
        shift_type_id: 0,
    });

    useEffect(() => {
        if (isOpen) {
            setLoading(true);
            Promise.all([
                employeeApi.list(),
                shiftTypeApi.list(),
            ]).then(([empRes, shiftRes]) => {
                setEmployees((empRes.data.data || []).filter((e: Employee) => e.IsActive));
                setShiftTypes(shiftRes.data);

                if (schedule) {
                    setFormData({
                        employee_id: schedule.EmployeeID,
                        shift_type_id: schedule.ShiftTypeID,
                    });
                } else {
                    setFormData({
                        employee_id: 0,
                        shift_type_id: shiftRes.data.length > 0 ? shiftRes.data[0].ID : 0,
                    });
                }
            }).catch(console.error).finally(() => setLoading(false));
        }
    }, [isOpen, schedule]);

    const handleSave = async () => {
        if (!formData.employee_id || !formData.shift_type_id) {
            alert('Lütfen personel ve nöbet tipi seçin.');
            return;
        }

        setSaving(true);
        try {
            // Adjust date to avoid timezone issues - set to noon
            const targetDate = new Date(date);
            targetDate.setHours(12, 0, 0, 0);

            const payload = {
                Date: targetDate.toISOString(),
                EmployeeID: formData.employee_id,
                ShiftTypeID: formData.shift_type_id,
                DepartmentID: employees.find(e => e.ID === formData.employee_id)?.DepartmentID || 0,
            };

            if (schedule) {
                await scheduleApi.update(schedule.ID, { ...schedule, ...payload });
            } else {
                await scheduleApi.create(payload);
            }
            onSuccess();
            onClose();
        } catch (err) {
            console.error(err);
            alert('Kaydetme başarısız: ' + err);
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = async () => {
        if (!schedule || !confirm('Bu nöbeti silmek istediğinize emin misiniz?')) return;

        setDeleting(true);
        try {
            await scheduleApi.delete(schedule.ID);
            onSuccess();
            onClose();
        } catch (err) {
            console.error(err);
            alert('Silme başarısız: ' + err);
        } finally {
            setDeleting(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-fade-in">
            <div className="glass-card w-full max-w-md p-6 relative animate-scale-in">
                <button
                    onClick={onClose}
                    className="absolute top-4 right-4 text-gray-400 hover:text-white"
                >
                    <X className="w-5 h-5" />
                </button>

                <h3 className="text-xl font-bold mb-1">
                    {schedule ? 'Nöbet Düzenle' : 'Yeni Nöbet Ekle'}
                </h3>
                <p className="text-sm text-gray-400 mb-6">
                    {date.toLocaleDateString('tr-TR', { day: 'numeric', month: 'long', year: 'numeric', weekday: 'long' })}
                </p>

                {loading ? (
                    <div className="flex justify-center py-8">
                        <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
                    </div>
                ) : (
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 flex items-center gap-2">
                                <User className="w-4 h-4" /> Personel
                            </label>
                            <select
                                className="glass-input w-full"
                                value={formData.employee_id}
                                onChange={(e) => setFormData({ ...formData, employee_id: Number(e.target.value) })}
                            >
                                <option value={0} disabled>Personel Seçin</option>
                                {employees.map(e => (
                                    <option key={e.ID} value={e.ID}>
                                        {e.FirstName} {e.LastName} ({e.Title?.Name})
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 flex items-center gap-2">
                                <Clock className="w-4 h-4" /> Nöbet Tipi
                            </label>
                            <div className="grid grid-cols-2 gap-2">
                                {shiftTypes.map(st => (
                                    <button
                                        key={st.ID}
                                        type="button"
                                        onClick={() => setFormData({ ...formData, shift_type_id: st.ID })}
                                        className={`flex flex-col items-start p-2 rounded-lg border text-sm transition-all ${formData.shift_type_id === st.ID
                                                ? 'bg-blue-500/20 border-blue-500/50 text-white'
                                                : 'bg-white/5 border-white/10 text-gray-400 hover:bg-white/10'
                                            }`}
                                    >
                                        <span className="font-semibold">{st.Name}</span>
                                        <span className="text-xs opacity-70">{st.StartTime}-{st.EndTime}</span>
                                    </button>
                                ))}
                            </div>
                        </div>

                        <div className="flex gap-3 mt-8 pt-4 border-t border-white/10">
                            {schedule && (
                                <button
                                    onClick={handleDelete}
                                    disabled={deleting || saving}
                                    className="btn-danger mr-auto"
                                >
                                    {deleting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
                                </button>
                            )}

                            <button onClick={onClose} className="btn-ghost" disabled={saving || deleting}>
                                İptal
                            </button>
                            <button
                                onClick={handleSave}
                                disabled={saving || deleting || !formData.employee_id || !formData.shift_type_id}
                                className="btn-primary flex-1 justify-center"
                            >
                                {saving ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <Save className="w-4 h-4 mr-2" />}
                                Kaydet
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ScheduleEditModal;
