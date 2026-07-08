// 📡 サーバー（Go）からWebPushの電波が届いた瞬間に発火するイベント
self.addEventListener('push', function(event) {
  if (!event.data) return;

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
    icon: '/icon-192.png',
    badge: '/icon-192.png',
    vibrate: [200, 100, 200],
    data: { url: '/' }
  };

  event.waitUntil(
    self.registration.showNotification(title, options)
  );
});

// 🖱️ 通知をタップしたときの動き
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
