import React, { useState } from 'react';
import { Wand2, CheckCircle2, ChevronRight, Loader2 } from 'lucide-react';
import { scheduleApi } from '../services/api';
import type { ScheduleRequest } from '../types';

interface Props {
    onNavigate?: (tab: string) => void;
}

const ScheduleWizard: React.FC<Props> = ({ onNavigate }) => {
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [params, setParams] = useState<ScheduleRequest>({
        month: new Date().getMonth() + 2 > 12 ? 1 : new Date().getMonth() + 2,
        year: new Date().getMonth() + 2 > 12 ? new Date().getFullYear() + 1 : new Date().getFullYear(),
        overtime_threshold: 45,
        overtime_multiplier: 1.5,
    });
    const [resultCount, setResultCount] = useState(0);

    const MONTHS_TR = [
        'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
        'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
    ];

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const res = await scheduleApi.generate(params);
            setResultCount(res.data?.length || 0);
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
                        Algoritmamız personelin haftalık çalışma saatlerini ve maliyeti optimize ederek en ideal listeyi oluşturur.
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
                            disabled={loading}
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
