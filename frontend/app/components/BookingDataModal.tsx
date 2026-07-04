import React from "react";

// 💡 ゴーサーバーから取得した「残り何組か」を親から受け取るための型定義
interface BookingDataModalProps {
  remainingGroups: number; // 残りの待機組数
  bookingNumber  : number; // 整理券番号
}

export default function BookingDataModal({ remainingGroups, bookingNumber }: BookingDataModalProps) {
  return (
    /* 外枠のプレミアムダークコンテナ */
    <div className="w-full mx-auto bg-[#1e1e22] rounded-2xl border border-zinc-700/10 ring-1 ring-cyan-400/25 pt-12 pb-14 px-8 flex flex-col items-center shadow-[0_30px_70px_rgba(0,0,0,0.6)]">
      
      {/* 1. 整理券番号セクション */}
      <div className="text-center w-full">
        <div className="mb-4">
          <span className="text-[10px] font-bold tracking-[0.25em] text-cyan-400/80 uppercase block mb-1">
            SUCCESSFULLY BOOKED
          </span>
          <span className="text-xs font-medium tracking-[0.15em] text-zinc-500 block">
            整理券番号
          </span>
        </div>
        <div className="flex items-baseline justify-center">
          <span className="text-8xl font-thin tracking-[-0.06em] text-cyan-300 drop-shadow-[0_0_15px_rgba(103,232,249,0.12)]">
            {bookingNumber}
          </span>
          <span className="text-xl font-medium text-zinc-400 ml-1 select-none">番</span>
        </div>
      </div>
      <div className="w-[calc(100%+4rem)] -mx-8 border-t border-dashed border-zinc-700/30 my-8"></div>

      {remainingGroups === 1 && (
        <div className="w-full bg-cyan-950/30 border border-cyan-500/30 rounded-xl p-4 mb-4 text-center">
          <span className="inline-block text-[10px] font-bold tracking-[0.2em] text-cyan-300 bg-cyan-900/50 rounded-md px-2 py-0.5 mb-1.5 select-none">
            お知らせ
          </span>
          <p className="text-xs font-semibold text-zinc-100 leading-relaxed">
            まもなくご案内いたします。<br />
            アトラクションの手前までお進みください。
          </p>
        </div>
      )}

      <div className="text-center bg-linear-to-b from-[#161619]/80 to-[#0e0e10]/95 border border-zinc-800/40 rounded-xl py-5 px-4 w-full shadow-inner">
        <span className="block text-[10px] font-bold tracking-[0.15em] text-zinc-400/80 uppercase mb-2">
          ご案内予定時刻
        </span>
        <div className="text-2xl font-light font-mono text-zinc-200 tracking-wide">
          14:35 <span className="text-xs font-medium text-zinc-500 ml-0.5">頃</span>
        </div>
      </div>

    </div>
  );
}
