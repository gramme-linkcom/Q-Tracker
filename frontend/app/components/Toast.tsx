import React, { useState, useEffect } from "react";

// 💡 親から受け取るトースト通知のデータ型
interface ToastProps {
  show: boolean;      // 表示中かどうか
  message: string;    // 表示する文字列
}

export default function Toast({ show, message }: ToastProps) {
  const [hasBeenShown, setHasBeenShown] = useState<boolean>(false);

  useEffect(() => {
    if (show) {
      setHasBeenShown(true);
    }
  }, [show]);

  if (!hasBeenShown) {
    return null;
  }

  return (
    <div className="toast toast-bottom toast-center z-50 pointer-events-none p-4 pb-8">
      <div 
        className={`alert bg-zinc-100 text-zinc-900 font-bold border-none shadow-2xl rounded-2xl p-4 flex items-center gap-3 transition-all duration-500 ease-out transform ${
          show 
            ? "opacity-100 translate-y-0" 
            : "opacity-0 translate-y-4"
        }`}
      >
        <svg xmlns="http://www.w3.org/2000/svg" className="stroke-emerald-600 shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span className="text-sm">{message}</span>
      </div>
    </div>
  );
}
