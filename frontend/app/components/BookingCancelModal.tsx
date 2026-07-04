import React from "react";
import { useState } from "react";

interface BookingCancelModalProps {
  isOpen: boolean;        // 開いているかどうか
  onClose: () => void;     // 閉じるボタンが押された時の関数
  onConfirm: () => void;   // 確定ボタンが押された時の関数
  bookingNumber: number;
}

export default function BookingCancelModal({ isOpen, onClose, onConfirm, bookingNumber }: BookingCancelModalProps) {
    return (
<div className={`modal modal-bottom sm:modal-middle transition-all duration-300 ${isOpen ? 'modal-open pointer-events-auto' : 'pointer-events-none'}`}>
      <div className="modal-box rounded-2xl border border-zinc-700/30 p-8 shadow-2xl bg-[#242428] max-w-sm mx-auto text-zinc-200">
        
        {/* ─── ここから下の中身は、今書いてくれているコードをそのまま入れる ─── */}
        <h3 className="font-bold text-lg text-center text-zinc-100 mb-2 flex items-center justify-center gap-2">
          <i className="bi bi-exclamation-triangle text-red-400 text-base"></i>
          予約のキャンセル
        </h3>
        
        <p className="text-xs text-center text-zinc-400 leading-relaxed mb-6">
          現在お持ちの整理券が削除されます。この操作は取り消せませんが、本当によろしいですか？
        </p>

        <div className="bg-[#1e1e22] rounded-xl border border-zinc-800 p-3 mb-6 text-center">
          <span className="text-[10px] font-bold text-zinc-500 block mb-0.5">キャンセルする整理券番号</span>
          <span className="text-xl font-light text-zinc-300 font-mono">No. <span className="font-normal text-cyan-300">{bookingNumber}</span></span>
        </div>
        
        <div className="flex gap-2">
          <button 
            type="button"
            onClick={onClose} // 💡 コメントアウトを外して、親から貰った onClose を繋ぐ
            className="btn btn-ghost rounded-xl flex-1 text-xs font-bold text-zinc-400 hover:bg-zinc-700/30"
          >
            戻る
          </button>
          <button 
            type="button"
            onClick={onConfirm} // 💡 コメントアウトを外して、親から貰った onConfirm を繋ぐ
            className="btn bg-red-950/40 text-red-400 border border-red-500/30 hover:bg-red-900/50 rounded-xl flex-1 text-xs font-bold tracking-wider"
          >
            キャンセルを確定
          </button>
        </div>
        {/* ─── ここまでの中身 ─── */}

      </div>

      {/* 3. 💡 ここもお好みで追加！ ポップアップの外の暗いエリアをタップした時も閉じるようにする仕掛け */}
      <div 
        onClick={onClose} 
        className="modal-backdrop bg-[#0f0f11]/70 backdrop-blur-xs cursor-pointer"
      ></div>

    </div>
    );
}
