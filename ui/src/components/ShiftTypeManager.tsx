import React, { useState, useEffect } from 'react';
import { Plus, Clock, Palette } from 'lucide-react';
import { shiftTypeApi } from '../services/api';

const ShiftTypeManager: React.FC = () => {
    const [shifts, setShifts] = useState<any[]>([]);
    const [showForm, setShowForm] = useState(false);
    const [formData, setFormData] = useState({
        Name: '',
        StartTime: '08:00',
        EndTime: '08:00',
        Color: '#3b82f6',
        Description: '',
    });

    const fetchShifts = async () => {
        try {
            const res = await shiftTypeApi.list();
            setShifts(res.data);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchShifts();
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await shiftTypeApi.create(formData);
            setShowForm(false);
            fetchShifts();
        } catch (err) {
            alert('Hata oluştu');
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Nöbet Tipleri</h2>
                <button
                    onClick={() => setShowForm(!showForm)}
                    className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg"
                >
                    <Plus className="w-4 h-4" />
                    Yeni Nöbet Tipi
                </button>
            </div>

            {showForm && (
                <form onSubmit={handleSubmit} className="bg-gray-700 p-6 rounded-lg grid grid-cols-2 gap-4">
                    <input
                        placeholder="Nöbet Adı (örn: 24 Saatlik)"
                        className="bg-gray-800 border-gray-600 rounded p-2"
                        value={formData.Name}
                        onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
                        required
                    />
                    <div className="flex items-center gap-2 bg-gray-800 p-1 rounded">
                        <Palette className="w-5 h-5 ml-2 text-gray-400" />
                        <input
                            type="color"
                            className="bg-transparent border-none w-10 h-8 p-0 cursor-pointer"
                            value={formData.Color}
                            onChange={(e) => setFormData({ ...formData, Color: e.target.value })}
                        />
                        <span className="text-sm text-gray-400">{formData.Color}</span>
                    </div>
                    <div className="flex flex-col gap-1">
                        <label className="text-xs text-gray-400">Başlangıç</label>
                        <input
                            type="time"
                            className="bg-gray-800 border-gray-600 rounded p-2"
                            value={formData.StartTime}
                            onChange={(e) => setFormData({ ...formData, StartTime: e.target.value })}
                            required
                        />
                    </div>
                    <div className="flex flex-col gap-1">
                        <label className="text-xs text-gray-400">Bitiş</label>
                        <input
                            type="time"
                            className="bg-gray-800 border-gray-600 rounded p-2"
                            value={formData.EndTime}
                            onChange={(e) => setFormData({ ...formData, EndTime: e.target.value })}
                            required
                        />
                    </div>
                    <textarea
                        placeholder="Açıklama"
                        className="bg-gray-800 border-gray-600 rounded p-2 col-span-2"
                        value={formData.Description}
                        onChange={(e) => setFormData({ ...formData, Description: e.target.value })}
                    />
                    <div className="col-span-2 flex justify-end gap-3 mt-2">
                        <button
                            type="button"
                            onClick={() => setShowForm(false)}
                            className="px-4 py-2 text-gray-400 hover:text-white"
                        >
                            İptal
                        </button>
                        <button type="submit" className="bg-green-600 hover:bg-green-700 px-6 py-2 rounded">
                            Kaydet
                        </button>
                    </div>
                </form>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {shifts.map((s) => (
                    <div key={s.ID} className="bg-gray-700 p-4 rounded-lg flex items-center gap-4 border-l-4" style={{ borderLeftColor: s.Color }}>
                        <div className="bg-gray-600 p-3 rounded-full">
                            <Clock className="w-6 h-6 text-gray-300" />
                        </div>
                        <div className="flex-1">
                            <div className="font-medium text-lg">{s.Name}</div>
                            <div className="text-gray-400 text-sm">
                                {s.StartTime} - {s.EndTime}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ShiftTypeManager;
