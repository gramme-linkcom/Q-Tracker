import React from "react";
import { useState } from "react";

// 💡 親（page.tsx）から受け取るデータの「型」を拡張します
interface BookingModalProps {
  isOpen: boolean;        // 開いているかどうか
  onClose: () => void;     // 閉じるボタンが押された時の関数
  onConfirm: (reservedTime: string) => void;   // 確定ボタンが押された時の関数
  isPending: boolean;
  serveStartTime: string;  // 稼働開始時間 (HH:MM)
  serveEndTime: string;    // 稼働終了時間 (HH:MM)
  slotInterval: number;    // 予約枠の粒度 (分)
  maxBookingsPerSlot: number; // 1枠あたりの最大予約数
  slotBookings: Record<string, number>; // 各時間枠の現在の予約数
  allowNoTimeSlot: boolean; // 時間指定なしでの発券を許可するかどうか
}

export default function BookingModal({ isOpen, onClose, onConfirm, isPending, serveStartTime, serveEndTime, slotInterval, maxBookingsPerSlot, slotBookings, allowNoTimeSlot }: BookingModalProps) {
  const [step, setStep] = useState("input");
  const [reservedTime, setReservedTime] = useState("");

  const slots = React.useMemo(() => {
    const list: string[] = [];
    try {
      const [startH, startM] = serveStartTime.split(":").map(Number);
      const [endH, endM] = serveEndTime.split(":").map(Number);
      
      let currentHour = startH;
      let currentMin = startM;
      
      const targetEndMinutes = endH * 60 + endM;
      
      while (true) {
        const nextMin = (currentMin + slotInterval) % 60;
        const nextHour = currentHour + Math.floor((currentMin + slotInterval) / 60);
        
        const currentTotal = currentHour * 60 + currentMin;
        const nextTotal = nextHour * 60 + nextMin;
        
        if (nextTotal > targetEndMinutes) {
          break;
        }
        
        const pad = (num: number) => String(num).padStart(2, '0');
        list.push(`${pad(currentHour)}:${pad(currentMin)} - ${pad(nextHour)}:${pad(nextMin)}`);
        
        currentHour = nextHour;
        currentMin = nextMin;
      }
    } catch (e) {
      console.error("Failed to generate slots:", e);
    }
    return list;
  }, [serveStartTime, serveEndTime, slotInterval]);

  const isSlotPast = (slot: string): boolean => {
    if (!slot) return false;
    try {
      const parts = slot.split(" - ");
      if (parts.length < 2) return false;
      const startTimeStr = parts[0];
      const [startHours, startMinutes] = startTimeStr.split(":").map(Number);
      
      const now = new Date();
      const currentHours = now.getHours();
      const currentMinutes = now.getMinutes();
      
      if (currentHours > startHours || (currentHours === startHours && currentMinutes >= startMinutes)) {
        return true;
      }
    } catch (e) {
      console.error(e);
    }
    return false;
  };

  const isAllSlotsPast = React.useMemo(() => {
    if (!serveEndTime) return false;
    try {
      const [endH, endM] = serveEndTime.split(":").map(Number);
      const now = new Date();
      const currentHours = now.getHours();
      const currentMinutes = now.getMinutes();
      
      if (currentHours > endH || (currentHours === endH && currentMinutes >= endM)) {
        return true;
      }
    } catch (e) {
      console.error(e);
    }
    return false;
  }, [serveEndTime]);

  const isConfirmDisabled = React.useMemo(() => {
    if (isPending) return true;
    if (reservedTime === "") {
      return isAllSlotsPast || !allowNoTimeSlot;
    }
    const past = isSlotPast(reservedTime);
    const bookedCount = slotBookings[reservedTime] || 0;
    const remaining = maxBookingsPerSlot - bookedCount;
    const isFull = remaining <= 0;
    return past || isFull;
  }, [reservedTime, isPending, isAllSlotsPast, allowNoTimeSlot, slotBookings, maxBookingsPerSlot]);

  const confirmButtonText = React.useMemo(() => {
    if (isPending) return "発行中...";
    if (reservedTime === "") {
      if (isAllSlotsPast) return "受付終了";
      if (!allowNoTimeSlot) return "時間指定必須";
      return "確定する";
    }
    if (isSlotPast(reservedTime)) {
      return "受付終了";
    }
    const bookedCount = slotBookings[reservedTime] || 0;
    const remaining = maxBookingsPerSlot - bookedCount;
    if (remaining <= 0) {
      return "満員";
    }
    return "確定する";
  }, [isPending, reservedTime, isAllSlotsPast, allowNoTimeSlot, slotBookings, maxBookingsPerSlot]);

  return (
    <div className={`modal modal-bottom sm:modal-middle transition-all duration-300 ${isOpen ? 'modal-open pointer-events-auto' : 'pointer-events-none'}`}>
      <div className="modal-box rounded-2xl border border-zinc-700/30 p-8 shadow-2xl bg-[#242428] max-w-sm mx-auto text-zinc-200">
        
        <h3 className="font-bold text-lg text-center text-zinc-100 mb-5">
          デジタル整理券の発行
        </h3>
        
        <div className="bg-[#1e1e22]/60 rounded-xl border border-zinc-800/40 p-4 mb-6">
          <span className="block text-[10px] font-bold tracking-[0.2em] text-zinc-500 uppercase mb-3 text-center">
            INFORMATION / 注意事項
          </span>
          
          <ul className="space-y-3 text-xs text-zinc-400 leading-relaxed">
            <li className="flex items-start gap-2.5">
              <i className="bi bi-clock-history text-cyan-400 text-sm mt-0.5 shrink-0"></i>
              <span>
                お呼び出し通知は、<strong className="text-zinc-100 font-semibold">直前のグループ入場時のご案内</strong>となります。通知が届きましたら、お早めに部屋の前までお越しください。
              </span>
            </li>
            <li className="flex items-start gap-2.5">
              <i className="bi bi-exclamation-circle text-zinc-500 text-sm mt-0.5 shrink-0"></i>
              <span>
                直前のグループの退場時（お呼び出しから<strong className="text-cyan-400 font-semibold">およそ5分以上</strong>）にいらっしゃらない場合は、予約を無効とし、次の方を先にご案内する場合がございます。
              </span>
            </li>
            <li className="flex items-start gap-2.5">
              <i className="bi bi-bell text-zinc-500 text-sm mt-0.5 shrink-0"></i>
              <span>
                通知を許可した場合、ブラウザを閉じても通知は届きますが、通信環境等により遅れる場合がございます。適宜画面を更新してご確認ください。
              </span>
            </li>
          </ul>
        </div>

        <div className="mb-6">
          <label className="block text-xs font-bold tracking-wider text-zinc-400 mb-2 select-none">
            希望のご案内時間帯
          </label>
          <select
            value={reservedTime}
            onChange={(e) => setReservedTime(e.target.value)}
            disabled={isPending}
            className="w-full bg-[#1e1e22] text-zinc-100 border border-zinc-800 rounded-xl px-4 py-3 text-sm focus:outline-hidden focus:ring-1 focus:ring-cyan-400 transition-all select-none"
          >
            <option value="" disabled={isAllSlotsPast || !allowNoTimeSlot}>
              時間指定なし（順番待ち）{isAllSlotsPast ? " (受付終了)" : !allowNoTimeSlot ? " (受付停止中)" : ""}
            </option>
            {slots.map((slot) => {
              const past = isSlotPast(slot);
              const bookedCount = slotBookings[slot] || 0;
              const remaining = maxBookingsPerSlot - bookedCount;
              const isFull = remaining <= 0;

              let label = slot;
              let isDisabled = past;

              if (past) {
                label = `${slot} (受付終了)`;
              } else if (isFull) {
                label = `${slot} (満員)`;
                isDisabled = true;
              } else {
                label = `${slot} (あと ${remaining} 枠)`;
              }

              return (
                <option key={slot} value={slot} disabled={isDisabled}>
                  {label}
                </option>
              );
            })}
          </select>
        </div>

        <p className="text-xs text-center text-zinc-400 leading-relaxed mb-6 px-2">
          順番が近づきましたら、スマートフォンへ通知が届きます。発行を確定してよろしいですか？
        </p>
        
        <div className="flex gap-2">
          {/* キャンセルボタン：通信中は押せないように制御 */}
          <button 
            type="button"
            onClick={onClose}
            disabled={isPending}
            className="btn btn-ghost rounded-xl flex-1 text-xs font-bold text-zinc-400 hover:bg-zinc-700/30 hover:text-zinc-200 disabled:opacity-40 select-none"
          >
            キャンセル
          </button>
          
          <button 
            type="button"
            onClick={() => onConfirm(reservedTime)}
            disabled={isConfirmDisabled}
            className={`btn rounded-xl flex-1 text-xs font-bold tracking-wider border-none select-none ${
              isConfirmDisabled 
                ? 'bg-zinc-800 text-zinc-500 cursor-not-allowed opacity-50' 
                : 'bg-zinc-100 text-zinc-900 hover:bg-zinc-200'
            }`}
          >
            {confirmButtonText}
          </button>
        </div>
      </div>
      
      {/* ポップアップ外の暗い背景部分：通信中はタップしても閉じないようにガード */}
      <div 
        onClick={() => { if (!isPending) onClose(); }} 
        className="modal-backdrop bg-[#0f0f11]/70 backdrop-blur-xs cursor-pointer"
      ></div>
    </div>
  );
}
