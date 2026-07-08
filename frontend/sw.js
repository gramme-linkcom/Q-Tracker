self.addEventListener('push', function(event) {
  if (!event.data) {
    console.warn('[Service Worker] Pushメッセージにデータが含まれていません。');
    return;
  }

  let title = "Q-Tracker 呼び出し";
  let body = "";

  try {
    const data = event.data.json();
    
    title = data.title || title;
    body = data.body || "";
  } catch (e) {
    body = event.data.text();
  }

  const options = {
    body: body,
    icon: '/icon-192.png',       // アプリのアイコン画像
    badge: '/icon-192.png',      // スマホの上のバーに出る小さなアイコン
    vibrate: [200, 100, 200],
    data: {
      url: '/'                   // 通知をタップしたときに開くURL
    }
  };

  event.waitUntil(
    self.registration.showNotification(title, options)
  );
});

self.addEventListener('notificationclick', function(event) {
  event.notification.close();

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then(function(clientList) {
      for (let i = 0; i < clientList.length; i++) {
        let client = clientList[i];
        if (client.url === '/' && 'focus' in client) {
          return client.focus();
        }
      }
      if (clients.openWindow) {
        return clients.openWindow(event.notification.data.url);
      }
    })
  );
});
