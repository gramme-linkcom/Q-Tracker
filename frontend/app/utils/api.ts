// 型定義
export interface QueueStatus {
  waitTime: number;
  timeRequired: number;
  waitingGroups: number;
  myAheadGroups : number;
  currentNumber: number;
  nextNumber: number;
  isActive: boolean;
  isBookingAvailable: boolean;
  isServiceAvailable: boolean;
  noticeMessage: string;
  infoMessage?: string;
  reservedTime?: string; // 指定された予約時間
  serveStartTime?: string; // 稼働開始時間 (HH:MM)
  serveEndTime?: string;   // 稼働終了時間 (HH:MM)
  slotInterval?: number;   // 予約枠の粒度 (分)
  maxBookingsPerSlot?: number; // 1枠あたりの最大予約数
  slotBookings?: Record<string, number>; // 各時間枠の現在の予約数
}

export async function getPublicVapidKey() {
  const res = await fetch('/api/vapid-public-key');
  if (!res.ok) {
      throw new Error('公開鍵の取得に失敗しました');
  }
  return (await res.text()).trim();
}

export async function getTicketExists(bookingNumber: number, uuid: string): Promise<boolean> {
  let url = `/api/exists-ticket?myNumber=${bookingNumber}&uuid=${uuid}`;
  const response = await fetch(url);
  if (!response.ok) throw new Error("データの取得に失敗しました");

  const data = await response.json();
  return data.isTicketAvailable;
}

// 現在の状況を取得する (GET)
export async function fetchQueueStatus(bookingNumber: number): Promise<QueueStatus> {
  let url = "/api/data";
  if (bookingNumber > 0) {
    url += `?myNumber=${bookingNumber}`;
  }

  const response = await fetch(url);
  
  if (!response.ok) throw new Error("データの取得に失敗しました");
  return response.json();
}

export async function bookTicket(pushToken: string = "", reservedTime: string = ""): Promise<{ bookingNumber: number; uuid: string; success: boolean }> {
  const response = await fetch("/api/booking", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ pushToken, reservedTime }),
  });
  if (!response.ok) {
    // もし上限エラーなどでステータスコード400が返ってきた場合、詳細なエラーメッセージを読み取る
    const errData = await response.json().catch(() => ({}));
    throw new Error(errData.error || "予約に失敗しました");
  }
  return response.json();
}

// 整理券をキャンセルする (POST)
export async function cancelTicket(bookingNumber: number): Promise<{ success: boolean }> {
  const response = await fetch("/api/booking/cancel", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ bookingNumber }),
  });
  if (!response.ok) throw new Error("キャンセルに失敗しました");
  return response.json();
}
