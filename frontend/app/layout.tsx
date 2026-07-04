import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "待ち時間",
  description: "現在の待ち時間を表示します。",

};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="ja"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <head>
      </head>
      <body className="min-h-full flex flex-col">
        {children}
        <footer className="bg-[#141416] text-center pb-3 text-[#242428]">
          <span className="flex w-auto items-center justify-center">
            <p className="my-4">developed by 9ramme</p>
            <p>©2026 9ramme.net</p>
          </span>
          <p>OSS Project <a href="https://github.com/gramme-linkcom/KitaF_Q-Tracker">Q-Tracker</a></p>
        </footer>
      </body>
    </html>
  );
}
