"use client";

import { useState, useEffect, useMemo, useRef } from "react";

import BookingModal from "./components/BookingModal";
import Toast from "./components/Toast";
import IosModal from "./components/IosModal";
import DetectIosBrowser from "./utils/DetectIosBrowser";
import BookingDataModal from "./components/BookingDataModal";
import BookingCancelModal from "./components/BookingCancelModal";
import { bookTicket, cancelTicket, fetchQueueStatus, getTicketExists, getPublicVapidKey } from "./utils/api";

export default function Home() {
  const REFRESH_INTERVAL_SEC = 30;
  // --- 状態管理 (SPAのコア) ---
  const [waitTime, setWaitTime] = useState<number>(0);
  const [waitingGroups, setWaitingGroups] = useState<number>(0); // 待機グループ数
  const [lastUpdateTime, setLastUpdateTime] = useState<string>("00:00");
  const [isBookingAvailable, setIsBookingAvailable] = useState<boolean>(false);
  const [bookingNumber, setBookingNumber] = useState<number>(0);
  const bookingNumberRef = useRef<number>(0);
  const [isServiceAvailable, setIsServiceAvailable] = useState<boolean>(false);
  const [isNotificationDenied, setIsNotificationDenied] = useState<boolean>(false);
  
  const [noticeMessage, setNoticeMessage] = useState<string>(""); 
  const [infoMessage, setInfoMessage] = useState<string>("");
  const [currentNumber, setCurrentNumber] = useState<number>(0);
  const [nextNumber, setNextNumber] = useState<number>(0);
  const [remainGroups, setRemainGroups] = useState<number>(0)
  const [timeRequired, setTimeRequired] = useState<number>(0);
  const [reservedTime, setReservedTime] = useState<string>("");
  const [serveStartTime, setServeStartTime] = useState<string>("09:00");
  const [serveEndTime, setServeEndTime] = useState<string>("17:00");
  const [slotInterval, setSlotInterval] = useState<number>(30);
  const [maxBookingsPerSlot, setMaxBookingsPerSlot] = useState<number>(5);
  const [slotBookings, setSlotBookings] = useState<Record<string, number>>({});
  const [allowNoTimeSlot, setAllowNoTimeSlot] = useState<boolean>(true);

  const [isBooked, setIsBooked] = useState<boolean>(false);
  
  // ★【追加】連打・多重送信を防止するための通信中フラグ
  const [isBookingInProgress, setIsBookingInProgress] = useState<boolean>(false);
  
  // ポップアップとトーストの開閉状態
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [showToast, setShowToast] = useState<boolean>(false);
  const [showIosModal, setShowIosModal] = useState<boolean>(false);
  const [isCancelModalOpen, setIsCancelModalOpen] = useState<boolean>(false);
  const [countdown, setCountdown] = useState<number>(REFRESH_INTERVAL_SEC);
  const [showBookingData, setShowBookingData] = useState<boolean>(false)

  const [pageTitle, setPageTitle] = useState<string>("Q-Tracker");
  const [roomName, setRoomName] = useState<string>("Room");

  useEffect(() => {
    bookingNumberRef.current = bookingNumber;
  }, [bookingNumber]);

  const updateTime = () => {
    const d = new Date();
    const hour = String(d.getHours()).padStart(2, '0');
    const minute = String(d.getMinutes()).padStart(2, '0');
    setLastUpdateTime(`${hour}:${minute}`);
  };

  const handleRefresh = async (forcedNumber?: number) => {
    updateTime();
    
    const localNumber = localStorage.getItem('booking_number');
    const localTicketUUID = localStorage.getItem('booking_uuid');
    if (localNumber != null && localTicketUUID != null) {
      const num = Number(localNumber);
      if (!await getTicketExists(num, localTicketUUID)){
        localStorage.removeItem('booking_number');
        localStorage.removeItem('booking_uuid');
        setBookingNumber(0);
        setReservedTime("");
        setShowBookingData(false);
        setIsBooked(false);
      }
    }

    setCountdown(REFRESH_INTERVAL_SEC);
    try {
      const targetNumber = forcedNumber !== undefined ? forcedNumber : bookingNumberRef.current;
      const data = await fetchQueueStatus(targetNumber);
      setTimeRequired(data.timeRequired);
      setWaitTime(data.waitTime);
      setWaitingGroups(data.waitingGroups);
      setCurrentNumber(data.currentNumber);
      setNextNumber(data.nextNumber);
      setNoticeMessage(data.noticeMessage);
      setIsBookingAvailable(data.isBookingAvailable);
      setIsServiceAvailable(data.isServiceAvailable);
      if (data.infoMessage !== undefined) setInfoMessage(data.infoMessage);
      setRemainGroups(data.myAheadGroups ?? 0);
      if (data.reservedTime !== undefined) {
        setReservedTime(data.reservedTime);
      }
      if (data.serveStartTime !== undefined) {
        setServeStartTime(data.serveStartTime);
      }
      if (data.serveEndTime !== undefined) {
        setServeEndTime(data.serveEndTime);
      }
      if (data.slotInterval !== undefined) {
        setSlotInterval(data.slotInterval);
      }
      if (data.maxBookingsPerSlot !== undefined) {
        setMaxBookingsPerSlot(data.maxBookingsPerSlot);
      }
      if (data.slotBookings !== undefined) {
        setSlotBookings(data.slotBookings);
      }
      if (data.allowNoTimeSlot !== undefined) {
        setAllowNoTimeSlot(data.allowNoTimeSlot);
      }
    } catch (error) {
      console.error(error);
    }
  };

  const BookingButtonText = useMemo(() => {
    if (showIosModal) {
      return "ホーム画面に追加してください";
    } else if (isBookingInProgress) {
      return "発行しています..."; // ★ 通信中であることがわかるようにテキストを変更
    } else if (isBooked) {
      return "予約をキャンセルする";
    } else if (isBookingAvailable && isServiceAvailable) {
      return "整理券を発行する";
    } else {
      return "予約停止中";
    }
  }, [isBooked, isBookingAvailable, isServiceAvailable, showIosModal, isBookingInProgress]);
  
  // 整理券発行
  const confirmBooking = async (reservedTime: string) => {
    if (isBookingInProgress) return;

    try {
      setIsBookingInProgress(true);

      let pushToken = "";
      try {
        // 1. サーバー（Go）からVAPID公開鍵をスマートに取得
        const vapidKey = await getPublicVapidKey();
        
        // 2. util.tsの関数を呼び出し、通知許可取得＆本日23:59期限付き端末IDを生成
        // （※読み込み元に import { requestNotificationPermission } from "./utils/util"; が必要です）
        const { requestNotificationPermission } = await import("./utils/util");
        pushToken = await requestNotificationPermission(vapidKey);
      } catch (pushErr) {
        // 通知登録フェーズでエラー（SafariのPWA化忘れや権限拒否など）が起きても、
        // 発券自体は手動と同じ扱いで継続できるように安全ガード（フォールバック）を敷きます
        console.warn("[WebPush] 通知トークンの取得をスキップして発券を続行します:", pushErr);
        pushToken = "";
      }

      // 3. 生成した端末ID（または空文字）を乗せてサーバーへ一撃ポスト！
      const data = await bookTicket(pushToken, reservedTime);
      
      await handleRefresh(data.bookingNumber);
      setBookingNumber(data.bookingNumber);
      localStorage.setItem('booking_number', `${data.bookingNumber}`);
      localStorage.setItem('booking_uuid', `${data.uuid}`);
      setShowBookingData(true);
      setIsBooked(true);
      setShowToast(true);
      setIsModalOpen(false);
    } catch (error: any) {
      alert(error.message || "整理券の発行に失敗しました。もう一度お試しください。");
      await handleRefresh();
    } finally {
      setIsBookingInProgress(false);
    }

    setIsNotificationDenied(false);
    if ("Notification" in window && Notification.permission === "denied") {
      setIsNotificationDenied(true);
    }
  };

// 整理券キャンセル
const confirmCancelBooking = async () => {
  try {
    await cancelTicket(bookingNumber);
    localStorage.removeItem('booking_number');
    localStorage.removeItem('booking_uuid');
    setIsBooked(false);
    setReservedTime("");
    setIsCancelModalOpen(false);
    setShowBookingData(false);
    handleRefresh();
  } catch (error) {
    alert("キャンセルの処理に失敗しました。");
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
    const loadData = async () => {
      const localNumber = localStorage.getItem('booking_number');
      const localTicketUUID = localStorage.getItem('booking_uuid');
      
      if (localNumber != null && localTicketUUID != null) {
        const num = Number(localNumber);
        if (await getTicketExists(num, localTicketUUID)){
          setBookingNumber(num);
          setShowBookingData(true);
          setIsBooked(true);
          await handleRefresh(num);
        } else {
          localStorage.removeItem('booking_number');
          localStorage.removeItem('booking_uuid');
          await handleRefresh()
        }

      } else {
        await handleRefresh();
      }
      updateTime();
    };

    loadData();
    setShowIosModal(DetectIosBrowser());

    if ('serviceWorker' in navigator) {
      window.addEventListener('load', function() {
        navigator.serviceWorker.register('/sw.js')
          .then(function(registration) {
            console.log('[ServiceWorker] 登録成功！スコープ: ', registration.scope);
          })
          .catch(function(error) {
            console.error('[ServiceWorker] 登録失敗... 原因:', error);
          });
      });
    };

    if (typeof window !== "undefined" && (window as any).__SERVER_CONFIG__) {
      if ((window as any).__SERVER_CONFIG__.pageTitle) setPageTitle((window as any).__SERVER_CONFIG__.pageTitle);
      if ((window as any).__SERVER_CONFIG__.roomName) setRoomName((window as any).__SERVER_CONFIG__.roomName);
    }

    if ("Notification" in window && Notification.permission === "denied") {
      setIsNotificationDenied(true);
    }
  }, []);

  return (
    <div>
      <header className="navbar fixed top-0 left-0 w-full h-14 bg-[#0e0e10]/80 backdrop-blur-md border-b border-zinc-800/50 flex items-center z-50">
        <div className="w-full mx-4 md:mx-auto max-w-md flex items-center justify-between">
          <p className="text-xl font-light tracking-[0.15em] text-zinc-100">
            {pageTitle}
          </p>
        </div>
      </header>
      <div className="min-h-screen w-full bg-[#141416] flex flex-col items-center justify-start py-20 px-6 font-sans antialiased text-zinc-300 transition-colors duration-300">
        <div className="w-full max-w-md flex flex-col gap-6">
          <IosModal show={showIosModal} />

          {infoMessage &&(
            <div className="w-full bg-cyan-950/20 border border-cyan-500/20 rounded-2xl p-5 shadow-lg flex items-start gap-3 animate-fade-in">
              <div className="flex flex-col gap-1">
                <span className="text-[10px] font-bold text-sm flex items-center tracking-[0.2em] text-cyan-400/90 uppercase">
                  <i className="bi bi-info-circle-fill text-cyan-400 mt-0.5 shrink-0 mr-2"></i>
                  NOTICE / 運営からのお知らせ
                </span>
                <p className="text-sm font-medium text-zinc-200 leading-relaxed">
                  {infoMessage}
                </p>
              </div>
            </div>
          )}
          {isNotificationDenied && isBooked && (
            <div className="w-full bg-amber-950/20 border border-amber-500/20 rounded-2xl p-5 shadow-lg flex items-start gap-3">
              <div className="flex flex-col gap-1 w-full">
                <span className="text-[10px] font-bold tracking-[0.2em] text-amber-400 bg-amber-950/50 rounded-md px-2 py-0.5 w-fit uppercase select-none flex items-center">
                  <i className="bi bi-exclamation-triangle-fill mr-2"></i>
                  WARNING / 通知設定オフ
                </span>
                <p className="text-xs font-medium text-zinc-200 leading-relaxed mt-1">
                  ブラウザの通知設定が「ブロック」されています。<br />
                  このままだと順番が来た際のお知らせ通知が届かないため、画面を開いたまま待ち時間をご確認ください。
                </p>
              </div>
            </div>
          )}

          <BookingDataModal
            currentNumber={currentNumber}
            nextNumber={nextNumber}
            bookingNumber={bookingNumber}
            timeRequired={timeRequired}
            remainGroups={remainGroups}
            show={showBookingData}
            reservedTime={reservedTime}
          />
          
          <div className="w-full bg-[#1e1e22] rounded-2xl border border-zinc-700/20 p-8 md:p-12 flex flex-col items-center shadow-2xl">
            
            {isServiceAvailable && (
              <div className="flex items-center gap-2 mb-2">
                <span className="inline-block w-2 h-2 rounded-full bg-emerald-400"></span>
                <span className="text-xs font-bold tracking-widest text-zinc-400 uppercase">OPEN</span>
              </div>
            )}
            {!isServiceAvailable && (
              <div className="flex items-center gap-2 mb-2">
                <span className="inline-block w-2 h-2 rounded-full bg-red-400 "></span>
                <span className="text-xs font-bold tracking-widest text-zinc-400 uppercase">CLOSE</span>
              </div>
            )}
            
            <h1 className="text-2xl font-bold tracking-tight text-zinc-100 mb-8">
              {roomName}
            </h1>

            {isServiceAvailable && (
              <div className="text-center mb-6">
                  <span className="block text-xs font-bold tracking-widest text-zinc-500 uppercase mb-1">
                    現在の待ち時間
                  </span>
                  <div className="flex items-baseline justify-center mb-2">
                    <span className="text-xl font-medium text-zinc-400 ml-2">約</span>
                    <span className="text-8xl font-light tracking-tighter text-cyan-300 drop-shadow-[0_0_12px_rgba(103,232,249,0.15)] transition-all duration-300">
                      {waitTime}
                    </span>
                    <span className="text-xl font-medium text-zinc-400 ml-2">分</span>
                  </div>

                  <div className="text-xs text-zinc-500 tracking-wider">
                    待機列: <span className="font-semibold text-zinc-300">{waitingGroups} 組</span>
                  </div>
              </div>
            )}
            {!isServiceAvailable && (
              <div className="text-center mb-6">
                  <div className="flex items-baseline justify-center mb-2">
                    <span className="text-4xl font-bold tracking-tighter text-zinc-400 drop-shadow-[0_0_12px_rgba(103,232,249,0.15)] transition-all duration-300">
                      受付時間外
                    </span>
                  </div>
              </div>
            )}

            <div className="text-xs text-zinc-500 font-mono tracking-wider mb-8">
              最終更新時刻 {lastUpdateTime}
            </div>

            <div className="text-[11px] text-cyan-400/80 font-mono tracking-widest uppercase mb-8 flex items-center justify-center gap-1.5">
              <i className="bi bi-arrow-clockwise"></i>
              <span>{countdown} 秒後に自動更新</span>
            </div>

            {isBookingAvailable && isServiceAvailable &&(
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
                // ★【修正】現在通信中（isBookingInProgress）の場合もボタンを無効化する
                disabled={!isBookingAvailable || !isServiceAvailable || showIosModal || isBookingInProgress}
                className={`btn btn-block rounded-xl h-12 text-sm font-bold tracking-wide transition-all ${
                  showIosModal || isBookingInProgress
                    ? 'bg-zinc-800/60 text-zinc-500 cursor-not-allowed opacity-40 border-zinc-700/10'
                    : isBooked
                      ? 'bg-red-950/40 text-red-400 border border-red-500/30 hover:bg-red-900/50 active:scale-[0.98]' 
                      : isBookingAvailable && isServiceAvailable
                        ? 'bg-zinc-100 text-zinc-900 hover:bg-zinc-200 active:scale-[0.98] border-zinc-700/10 shadow-sm'
                        : 'bg-zinc-800/60 text-zinc-500 cursor-not-allowed opacity-40 border-zinc-700/10'
                }`}
              >
                {BookingButtonText}
              </button>

              <button 
                onClick={() => handleRefresh()}
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
          isPending={isBookingInProgress}
          serveStartTime={serveStartTime}
          serveEndTime={serveEndTime}
          slotInterval={slotInterval}
          maxBookingsPerSlot={maxBookingsPerSlot}
          slotBookings={slotBookings}
          allowNoTimeSlot={allowNoTimeSlot}
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
