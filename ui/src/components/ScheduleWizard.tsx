import React, { useState } from 'react';
import { Wand2, CheckCircle2, ChevronRight, Loader2 } from 'lucide-react';
import { scheduleApi } from '../services/api';

const ScheduleWizard: React.FC = () => {
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [params, setParams] = useState({
        month: new Date().getMonth() + 2, // Next month
        year: new Date().getFullYear(),
        overtime_threshold: 45,
        overtime_multiplier: 1.5,
    });
    const [result, setResult] = useState<any[] | null>(null);

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const res = await scheduleApi.generate(params);
            setResult(res.data);
            setStep(3);
        } catch (err) {
            alert('Oluşturma hatası: ' + err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto py-8">
            {/* Stepper */}
            <div className="mb-12 flex justify-between items-center relative">
                <div className="absolute top-1/2 left-0 w-full h-0.5 bg-gray-700 -z-10"></div>
                {[1, 2, 3].map((num) => (
                    <div
                        key={num}
                        className={`w-10 h-10 rounded-full flex items-center justify-center font-bold border-2 transition-colors ${step >= num ? 'bg-blue-600 border-blue-600 text-white' : 'bg-gray-800 border-gray-600 text-gray-500'
                            }`}
                    >
                        {step > num ? <CheckCircle2 className="w-6 h-6" /> : num}
                    </div>
                ))}
            </div>

            {step === 1 && (
                <div className="space-y-6 text-center">
                    <Wand2 className="w-16 h-16 mx-auto text-blue-400 opacity-50" />
                    <h3 className="text-2xl font-bold">Akıllı Nöbet Planlayıcı</h3>
                    <p className="text-gray-400">
                        Algoritmamız personelin haftalık çalışma saatlerini ve maliyeti optimize ederek en ideal listeyi oluşturur.
                    </p>
                    <button
                        onClick={() => setStep(2)}
                        className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 px-8 py-3 rounded-xl mx-auto text-lg font-semibold transition-all hover:scale-105"
                    >
                        Başlayalım <ChevronRight className="w-5 h-5" />
                    </button>
                </div>
            )}

            {step === 2 && (
                <div className="space-y-6 bg-gray-700 p-8 rounded-2xl shadow-2xl animate-in fade-in slide-in-from-bottom-5">
                    <h3 className="text-xl font-bold mb-4">Planlama Parametreleri</h3>
                    <div className="grid grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400">Yıl</label>
                            <input
                                type="number"
                                className="w-full bg-gray-800 border-none rounded-lg p-3"
                                value={params.year}
                                onChange={(e) => setParams({ ...params, year: Number(e.target.value) })}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400">Ay</label>
                            <select
                                className="w-full bg-gray-800 border-none rounded-lg p-3"
                                value={params.month}
                                onChange={(e) => setParams({ ...params, month: Number(e.target.value) })}
                            >
                                {Array.from({ length: 12 }, (_, i) => (
                                    <option key={i + 1} value={i + 1}>{i + 1}</option>
                                ))}
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400">Ek Mesai Eşiği (Saat)</label>
                            <input
                                type="number"
                                className="w-full bg-gray-800 border-none rounded-lg p-3"
                                value={params.overtime_threshold}
                                onChange={(e) => setParams({ ...params, overtime_threshold: Number(e.target.value) })}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm text-gray-400">Mesai Çarpanı (x)</label>
                            <input
                                type="number"
                                step="0.1"
                                className="w-full bg-gray-800 border-none rounded-lg p-3"
                                value={params.overtime_multiplier}
                                onChange={(e) => setParams({ ...params, overtime_multiplier: Number(e.target.value) })}
                            />
                        </div>
                    </div>
                    <button
                        onClick={handleGenerate}
                        disabled={loading}
                        className="w-full flex items-center justify-center gap-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 py-4 rounded-xl text-lg font-bold mt-4"
                    >
                        {loading ? <Loader2 className="w-6 h-6 animate-spin" /> : 'Sihirbazı Çalıştır'}
                    </button>
                </div>
            )}

            {step === 3 && (
                <div className="space-y-6 text-center animate-in zoom-in-95">
                    <div className="w-20 h-20 bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4 border-2 border-green-500">
                        <CheckCircle2 className="w-12 h-12 text-green-500" />
                    </div>
                    <h3 className="text-2xl font-bold text-white">Çizelge Başarıyla Oluşturuldu!</h3>
                    <p className="text-gray-400">
                        {result?.length} adet nöbet ataması optimize edilerek veritabanına kaydedildi.
                    </p>
                    <div className="flex gap-4 justify-center mt-8">
                        <button
                            onClick={() => setStep(1)}
                            className="px-6 py-2 border border-gray-600 rounded-lg hover:bg-gray-700"
                        >
                            Yeni Plan
                        </button>
                        <button className="px-6 py-2 bg-blue-600 rounded-lg hover:bg-blue-700">
                            Takvime Git
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ScheduleWizard;
