export default function IosModal({ show }: {show: boolean} ) {
  if (!show) {
    return null
  }

  return (
    <div className="w-full bg-[#141416]/50 rounded-xl border border-zinc-850 p-4 space-y-3 text-xs md:text-sm text-zinc-400 leading-relaxed">
    <span className="flex">
        <i className="bi bi-exclamation-circle-fill mr-2 text-warning"></i><h2 className="font-bold text-zinc-300">Apple製デバイスのブラウザからご利用ですか?</h2>
    </span>
    <p className="text-zinc-300 font-medium">
        Apple製デバイス（iPhone/iPad）のブラウザ制限により、このままでは呼び出し通知が届きません。
    </p>
    
    <div className="flex flex-col gap-2 pt-2 border-t border-zinc-800/40">
        <div className="flex items-center gap-2">
        <span className="badge bg-zinc-800 text-zinc-500 border-none font-mono text-[10px] w-5 h-5 p-0">1</span>
        <span>
            共有ボタン
            <span className="inline-flex items-center justify-center bg-zinc-800 border border-zinc-700 rounded-md p-1 mx-1 text-zinc-200">
            <i className="bi bi-box-arrow-up text-xs"></i>
            </span>
            をタップします。
        </span>
        </div>
        
        <div className="flex items-center gap-2">
        <span className="badge bg-zinc-800 text-zinc-500 border-none font-mono text-[10px] w-5 h-5 p-0">2</span>
        <span>
            メニュー内の<strong className="text-cyan-400 font-semibold">「ホーム画面に追加」</strong>を選択します。
        </span>
        </div>
    </div>
</div>
  );
}
