"use client";

import { useState, useEffect } from "react";
import Image from "next/image";

export default function Home() {

  const [lastUpdateTime, lastUpdateTime_updater] = useState("00:00");
  
  const GetTime = () => {
    let d = new Date();
    let [hour, minute] = [d.getHours(), d.getMinutes()]
    lastUpdateTime_updater(`${hour}:${minute}`);
  }

  useEffect(() => {
    GetTime();
  }, []);

  return (
    <div className="w-full max-w-sm md:max-w-2xl mx-auto flex flex-col gap-4 mt-8 md:mt-12 px-4">

        <div className="card w-full bg-base-100 shadow-2xl border-t-4 border-primary">
            <div className="card-body items-center text-center p-6 md:p-12">
                <div className="badge badge-secondary badge-lg font-bold mb-2 shadow-sm md:p-4 md:text-lg">営業中</div>
                
                <h2 className="card-title text-2xl md:text-4xl mb-4 font-bold">H301</h2>

                <div className="stat p-0 my-4 md:my-8">
                    <div className="stat-title text-base md:text-xl font-bold">現在の待ち時間</div>
                    <div className="stat-value text-7xl md:text-[9rem] text-primary my-2 drop-shadow-sm transition-all">
                        15<span className="text-3xl md:text-5xl text-base-content"> 分</span>
                    </div>
                    <div className="stat-desc text-sm md:text-base mt-2 font-medium text-base-content/70">最終更新: {lastUpdateTime}</div>
                </div>

                <div className="card-actions mt-6 w-full md:w-2/3 mx-auto">
                    <button 
                        className="btn btn-primary btn-block text-lg md:text-xl md:h-16 font-bold shadow-md hover:scale-[1.02] transition-transform" 
                    >
                        最新の状態に更新
                    </button>
                </div>
            </div>
        </div>

        <div className="alert alert-info shadow-sm mt-4 text-xs md:text-sm font-bold justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-5 h-5">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span>実際の待ち時間と多少異なる場合があります。</span>
        </div>

    </div>
  );
}
