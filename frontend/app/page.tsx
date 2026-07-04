"use client";

import { useState, useEffect } from "react";

// 💡 先ほど別ファイルで作ったコンポーネント部品をインポートする！
import BookingModal from "./components/BookingModal";
import Toast from "./components/Toast";
import IosModal from "./components/IosModal";
import DetectIosBrowser from "./components/DetectIosBrowser";
import BookingDataModal from "./components/BookingDataModal";
import BookingCancelModal from "./components/BookingCancelModal";

export default function Home() {
  const REFRESH_INTERVAL_SEC = 30;
  // --- 状態管理 (SPAのコア) ---
  const [waitTime, setWaitTime] = useState<number>(15);
  const [waitingGroups, setWaitingGroups] = useState<number>(3); // 待機グループ数
  const [lastUpdateTime, setLastUpdateTime] = useState<string>("00:00");
  const [isBookingAvailable, setIsBookingAvailable] = useState<boolean>(true);
  
  const [noticeMessage, setNoticeMessage] = useState<string>(""); 
  const [infoMessage, setInfoMessage] = useState<string>("");

  const [BookingButtonText, setBookingButtonText] = useState<string>("予約停止中") 
  const [isBooked, setIsBooked] = useState<boolean>(false);
  
  // ポップアップとトーストの開閉状態
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [showToast, setShowToast] = useState<boolean>(false);
  const [showIosModal, setShowIosModal] = useState<boolean>(false);
  const [isCancelModalOpen, setIsCancelModalOpen] = useState<boolean>(false);
  const [countdown, setCountdown] = useState<number>(REFRESH_INTERVAL_SEC);

  const [pageTitle, setPageTitle] = useState<string>(() => {
    if (typeof window !== "undefined" && (window as any).__SERVER_CONFIG__) {
      return (window as any).__SERVER_CONFIG__.pageTitle;
    }
    return "Q-Tracker"; // フォールバック
  });

  const [roomName, setRoomName] = useState<string>(() => {
    if (typeof window !== "undefined" && (window as any).__SERVER_CONFIG__) {
      return (window as any).__SERVER_CONFIG__.roomName;
    }
    return "Room"; // フォールバック
  });

  const bookingNumber = 1;

  const updateTime = () => {
    const d = new Date();
    const hour = String(d.getHours()).padStart(2, '0');
    const minute = String(d.getMinutes()).padStart(2, '0');
    setLastUpdateTime(`${hour}:${minute}`);
  };

  const handleRefresh = () => {
    updateTime();
    setCountdown(REFRESH_INTERVAL_SEC);
  };

  // 整理券発行の確定処理
  const confirmBooking = () => {
    setIsModalOpen(false);
    setShowToast(true);
  };

  // 予約キャンセルの確定処理
  const confirmCancelBooking = () => {
    setIsCancelModalOpen(false)
  }

  const setBookingButtonData = () => {
    if (DetectIosBrowser()) {
      setBookingButtonText("ホーム画面に追加してください");
      setIsBookingAvailable(false);
    } else if (isBooked) {
      setBookingButtonText("予約をキャンセルする");
    } else if (isBookingAvailable) {
      setBookingButtonText("整理券を発行する");
    } else {
      setBookingButtonText("予約停止中");
    }
  };

  // トーストがONになったら3秒後に自動で消す
  useEffect(() => {
    if (showToast) {
      const timer = setTimeout(() => {
        setShowToast(false);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [showToast]);

  useEffect(() => {
    const timer = setInterval(() => {
      setCountdown((prevCount) => {
        if (prevCount <= 1) {
          handleRefresh();
          return REFRESH_INTERVAL_SEC;
        }
        return prevCount - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    updateTime();
    setShowIosModal(DetectIosBrowser())
    setBookingButtonData()
  }, []);

  return (
    <div>
      <header className="navbar fixed top-0 left-0 w-full h-14 bg-[#0e0e10]/80 backdrop-blur-md border-b border-zinc-800/50 flex items-center z-50">
        <div className="w-full max-w-md mx-auto px-6 flex items-center justify-between">
          <p className="text-xl font-light tracking-[0.15em] text-zinc-100">
            {pageTitle}
          </p>
        </div>
      </header>
      <div className="min-h-screen w-full bg-[#141416] flex flex-col items-center justify-start py-20 px-6 font-sans antialiased text-zinc-300 transition-colors duration-300">
        <div className="w-full max-w-md flex flex-col gap-6">
          <IosModal show={showIosModal} />

          {infoMessage && (
            <div className="w-full bg-cyan-950/20 border border-cyan-500/20 rounded-2xl p-5 shadow-lg flex items-start gap-3 animate-fade-in">
              <i className="bi bi-info-circle-fill text-cyan-400 text-sm mt-0.5 shrink-0"></i>
              <div className="flex flex-col gap-1">
                <span className="text-[10px] font-bold tracking-[0.2em] text-cyan-400/90 uppercase">
                  NOTICE / 運営からのお知らせ
                </span>
                <p className="text-xs font-medium text-zinc-200 leading-relaxed">
                  {infoMessage}
                </p>
              </div>
            </div>
          )}

          <BookingDataModal
            remainingGroups={1}
            bookingNumber={bookingNumber}
          />
          
          <div className="w-full bg-[#1e1e22] rounded-2xl border border-zinc-700/20 p-8 md:p-12 flex flex-col items-center shadow-2xl">
            
            <div className="flex items-center gap-2 mb-2">
              <span className="inline-block w-2 h-2 rounded-full bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.4)] animate-pulse"></span>
              <span className="text-xs font-bold tracking-widest text-zinc-400 uppercase">OPEN</span>
            </div>
            
            <h1 className="text-2xl font-bold tracking-tight text-zinc-100 mb-8">
              {roomName}
            </h1>

            <div className="text-center mb-6">
              <span className="block text-xs font-bold tracking-widest text-zinc-500 uppercase mb-1">
                現在の待ち時間
              </span>
              <div className="flex items-baseline justify-center mb-2">
                <span className="text-8xl font-light tracking-tighter text-cyan-300 drop-shadow-[0_0_12px_rgba(103,232,249,0.15)] transition-all duration-300">
                  {waitTime}
                </span>
                <span className="text-xl font-medium text-zinc-400 ml-2">分</span>
              </div>

              <div className="text-xs text-zinc-500 tracking-wider">
                待機列: <span className="font-semibold text-zinc-300">{waitingGroups} 組</span>
              </div>
            </div>

            <div className="text-xs text-zinc-500 font-mono tracking-wider mb-8">
              最終更新時刻 {lastUpdateTime}
            </div>

            <div className="text-[11px] text-cyan-400/80 font-mono tracking-widest uppercase mb-8 flex items-center justify-center gap-1.5">
              <i className="bi bi-arrow-clockwise"></i>
              <span>{countdown} 秒後に自動更新</span>
            </div>

            {noticeMessage && (
              <div className="w-full border-t border-b border-zinc-700/20 py-4 text-center mb-8">
                <p className="text-sm font-medium text-zinc-300 leading-relaxed">
                  {noticeMessage}
                </p>
              </div>
            )}

            <div className="w-full flex flex-col gap-3">
              <button 
                onClick={() => {
                    if (!isBooked) {
                      setIsModalOpen(true)
                    } else {
                      setIsCancelModalOpen(true)
                    }
                  }}
                disabled={!isBookingAvailable || showIosModal}
                className={`btn btn-block rounded-xl h-12 text-sm font-bold tracking-wide transition-all ${
                  showIosModal
                    ? 'bg-zinc-800/60 text-zinc-500 cursor-not-allowed opacity-40 border-zinc-700/10'
                    : isBooked
                      ? 'bg-red-950/40 text-red-400 border border-red-500/30 hover:bg-red-900/50 active:scale-[0.98]' 
                      : isBookingAvailable
                        ? 'bg-zinc-100 text-zinc-900 hover:bg-zinc-200 active:scale-[0.98] border-zinc-700/10 shadow-sm'
                        : 'bg-zinc-800/60 text-zinc-500 cursor-not-allowed opacity-40 border-zinc-700/10'
                }`}
              >
                {BookingButtonText}
              </button>

              <button 
                onClick={handleRefresh}
                className="btn btn-ghost btn-block rounded-xl h-12 text-xs font-bold tracking-widest text-zinc-400 hover:text-zinc-200 hover:bg-zinc-700/20"
              >
                最新の状態にする
              </button>
            </div>

          </div>

          <p className="text-[11px] text-center text-zinc-500 font-medium leading-relaxed px-4">
            ※ 実際の混雑状況によって、ご案内までの時間は多少前後する場合があります。
          </p>

        </div>

        <BookingModal 
          isOpen={isModalOpen}                      
          onClose={() => setIsModalOpen(false)}     
          onConfirm={confirmBooking}                
        />
        <BookingCancelModal
          isOpen={isCancelModalOpen}
          onClose={() => setIsCancelModalOpen(false)}
          onConfirm={confirmCancelBooking}
          bookingNumber={bookingNumber}
        />

        <Toast 
          show={showToast} 
          message="整理券を発行しました！" 
        />

      </div>
    </div>
  );
}
