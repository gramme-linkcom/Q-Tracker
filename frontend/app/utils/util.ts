// util.ts

export function urlBase64ToUint8Array(base64String: string) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding).replace(/\-/g, '+').replace(/_/g, '/');
    const rawData = window.atob(base64);
    return new Uint8Array([...rawData].map(char => char.charCodeAt(0)));
}

function getTodayEndTimestamp(): number {
    const now = new Date();
    const endOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 999);
    return endOfToday.getTime();
}

export async function requestNotificationPermission(publicVapidKey: string): Promise<string> {
    // 1. 環境ガード
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
        console.warn('[WebPush] 未対応ブラウザです。');
        return "";
    }

    try {
        // 2. ユーザー操作の直後であるこのタイミングで通知許可を取る
        const permission = await Notification.requestPermission();
        if (permission !== 'granted') {
            console.warn('[WebPush] 拒否されました。');
            return "";
        }

        const reg = await navigator.serviceWorker.register('/sw.js');
        const registration = await navigator.serviceWorker.ready;

        const oldSub = await registration.pushManager.getSubscription();
        if (oldSub) {
            await oldSub.unsubscribe();
        }

        // 5. 新しい公開鍵でカギを生成
        const applicationServerKey = urlBase64ToUint8Array(publicVapidKey.trim());
        const subscription = await registration.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: applicationServerKey
        });

        // 6. 本日の23:59期限をカスタム合体させる
        const subscriptionRaw = subscription.toJSON();
        const subscriptionWithExpiry = {
            ...subscriptionRaw,
            expirationTime: getTodayEndTimestamp()
        };

        console.log('[WebPush] トークン取得成功！');
        return JSON.stringify(subscriptionWithExpiry);

    } catch (error) {
        console.error('[WebPush] エラー発生:', error);
        return "";
    }
}
