export function urlBase64ToUint8Array(base64String: string) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding).replace(/\-/g, '+').replace(/_/g, '/');
    const rawData = window.atob(base64);
    return new Uint8Array([...rawData].map(char => char.charCodeAt(0)));
}

function getTodayEndTimestamp(): number {
    const now = new Date();
    // 今日の 23時 59分 59秒 999ミリ秒 にセット
    const endOfToday = new Date(
        now.getFullYear(),
        now.getMonth(),
        now.getDate(),
        23,
        59,
        59,
        999
    );
    return endOfToday.getTime();
}

export async function requestNotificationPermission(publicVapidKey: string): Promise<string> {
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
        console.warn('[WebPush] このブラウザはPush通知をサポートしていません。');
        return "";
    }

    try {
        // 2. スマホの画面に「通知を許可しますか？」のポップアップを出す
        const permission = await Notification.requestPermission();
        if (permission !== 'granted') {
            console.warn('[WebPush] 通知権限がユーザーによって拒否されました。');
            return "";
        }

        // 3. バックグラウンドで動く Service Worker が準備万端になるまで待つ
        const registration = await navigator.serviceWorker.ready;

        // 4. すでに登録済みの通知用カギがあるか一回覗いてみる
        let subscription = await registration.pushManager.getSubscription();

        // 5. 【重要】今日の23:59で期限が切れるように、オプションを組んで新規登録（または上書き）する
        const applicationServerKey = urlBase64ToUint8Array(publicVapidKey);
        
        // 今日の23:59のタイムスタンプを取得
        const expirationTime = getTodayEndTimestamp();

        subscription = await registration.pushManager.subscribe({
            userVisibleOnly: true, // ユーザーに必ず見える通知を出すというブラウザへの誓約
            applicationServerKey: applicationServerKey
        });

        // 6. 取得した端末ID（構造体）に、フロント側でも有効期限をカスタム付与したオブジェクトを作る
        const subscriptionRaw = subscription.toJSON();
        const subscriptionWithExpiry = {
            ...subscriptionRaw,
            // サーバー側でも「今日の23:59」を検知してクリーンアップできるよう、期限をパケットに含める
            expirationTime: expirationTime 
        };

        console.log('[WebPush] 端末ID（トークン）の取得に成功しました。期限: 本日23:59まで');
        // JSON文字列にして、そのまま api/booking の pushToken に突っ込める形で返す
        return JSON.stringify(subscriptionWithExpiry);

    } catch (error) {
        console.error('[WebPush] 通知登録の処理中にエラーが発生しました:', error);
        return "";
    }
}
