import { MetadataRoute } from 'next'
export const dynamic = "force-static";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: 'デジタル整理券システム',
    short_name: 'Q-Tracker',
    description: '学校祭アトラクション順番待ちアプリ',
    start_url: '/',
    display: 'standalone', // 👈 これを指定すると、URLバーが消えて「完全なアプリ風」になります
    background_color: '#141416', // あなたのスペースグレー背景
    theme_color: '#1e1e22',
    icons: [
      { src: '/icon-192.png', sizes: '192x192', type: 'image/png' },
      { src: '/icon-512.png', sizes: '512x512', type: 'image/png' },
    ],
  }
}
