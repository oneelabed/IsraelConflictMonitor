"use client"

import { useState, useEffect } from "react";
import Link from "next/link"
import Image from "next/image"

export default function Navbar() {
  const [apiKey, setApiKey] = useState<string | null>(null);

  const checkAuth = () => {
    const key = localStorage.getItem("api_key");
    setApiKey(key);
  };

  useEffect(() => {
    checkAuth();

    window.addEventListener("storage-update", checkAuth);

    window.addEventListener("storage", checkAuth);

    return () => {
      window.removeEventListener("storage-update", checkAuth);
      window.removeEventListener("storage", checkAuth);
    };
  }, []);

  return (
    <nav className="sticky top-0 z-100 w-full border-b bg-white/100 backdrop-blur supports-[backdrop-filter]:bg-white/100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          
          {/* Left Side: Logo and Title */}
          <div className="flex items-center gap-3">
            <Link href="/" className="flex items-center gap-2">
              <div>
                <Image 
                    src="/ICMLogo.png" 
                    alt="ICM Logo" 
                    width={90} 
                    height={90}
                    className="object-contain" 
                />
              </div>
              <span className="hidden sm:inline-block font-bold text-xl tracking-tight text-slate-900">
                Israel Conflict Monitor
              </span>
            </Link>
          </div>

          {/* Right Side: Auth Links */}
          <div className="flex items-center gap-6">
            {apiKey ? (
              <>
                <Link 
                  href="/posts" 
                  className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
                >My Feed</Link>
                <Link 
                  href="/feeds" 
                  className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
                >Discover Feeds</Link>
                <button 
                  onClick={() => {
                    localStorage.removeItem("api_key");
                    localStorage.removeItem("username");
                    window.location.href = "/"; // Forces a hard refresh to clear the state
                  }}
                  className="text-sm font-medium text-red-500 hover:text-red-600 transition-colors"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link 
                  href="/login" 
                  className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
                >
                  Login
                </Link>
                <Link 
                  href="/signup" 
                  className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors"
                >
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}