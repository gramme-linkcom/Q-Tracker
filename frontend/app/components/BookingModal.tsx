import React from "react";
import { useState } from "react";

// 💡 親（page.tsx）から受け取るデータの「型」を定義します
interface BookingModalProps {
  isOpen: boolean;        // 開いているかどうか
  onClose: () => void;     // 閉じるボタンが押された時の関数
  onConfirm: () => void;   // 確定ボタンが押された時の関数
}

export default function BookingModal({ isOpen, onClose, onConfirm }: BookingModalProps) {
  const [step, setStep] = useState("input");

  return (
    <div className={`modal modal-bottom sm:modal-middle transition-all duration-300 ${isOpen ? 'modal-open pointer-events-auto' : 'pointer-events-none'}`}>
      <div className="modal-box rounded-2xl border border-zinc-700/30 p-8 shadow-2xl bg-[#242428] max-w-sm mx-auto text-zinc-200">
        
        <h3 className="font-bold text-lg text-center text-zinc-100 mb-5">
          デジタル整理券の発行
        </h3>
        
        {/* 💡 劇的に見やすくなった注意事項セクション */}
        <div className="bg-[#1e1e22]/60 rounded-xl border border-zinc-800/40 p-4 mb-6">
          <span className="block text-[10px] font-bold tracking-[0.2em] text-zinc-500 uppercase mb-3 text-center">
            INFORMATION / 注意事項
          </span>
          
          <ul className="space-y-3 text-xs text-zinc-400 leading-relaxed">
            <li className="flex items-start gap-2.5">
              <i className="bi bi-clock-history text-cyan-400 text-sm mt-0.5 shrink-0"></i>
              <span>
                お呼び出し通知は、<strong className="text-zinc-100 font-semibold">直前のグループの入場時のご案内</strong>となります。通知が届きましたら、お早めに部屋の前までお越しください。
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
                ブラウザを閉じても通知は届きますが、通信環境等により遅れる場合がございます。適宜画面を更新してご確認ください。
              </span>
            </li>
          </ul>
        </div>

        <p className="text-xs text-center text-zinc-400 leading-relaxed mb-6 px-2">
          順番が近づきましたら、スマートフォンへ通知が届きます。発行を確定してよろしいですか？
        </p>
        
        <div className="flex gap-2">
          {/* キャンセルボタン */}
          <button 
            type="button"
            onClick={onClose}
            className="btn btn-ghost rounded-xl flex-1 text-xs font-bold text-zinc-400 hover:bg-zinc-700/30 hover:text-zinc-200"
          >
            キャンセル
          </button>
          {/* 確定ボタン */}
          <button 
            type="button"
            onClick={onConfirm}
            className="btn bg-zinc-100 text-zinc-900 hover:bg-zinc-200 border-none rounded-xl flex-1 text-xs font-bold tracking-wider"
          >
            確定する
          </button>
        </div>
      </div>
      
      {/* ポップアップ外の暗い背景部分 */}
      <div 
        onClick={onClose} 
        className="modal-backdrop bg-[#0f0f11]/70 backdrop-blur-xs cursor-pointer"
      ></div>
    </div>
  );
}
