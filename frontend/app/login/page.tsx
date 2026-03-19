"use client"

import { useState } from "react";
import Link from "next/link";
import { AlertCircle, Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login } from "@/lib/actions/login";
import { useRouter } from "next/navigation";
import Cookies from 'js-cookie';

export default function SignUp() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const userData = await login(username.toLowerCase(), password);
      
      // 'expires: 1' means it lasts for 1 day
      Cookies.set("api_key", userData.api_key, { expires: 1, path: '/' });
      Cookies.set("role", userData.role, { expires: 1, path: '/' });
      Cookies.set("username", userData.name, { expires: 1, path: '/' });

      console.log("Logged in successfully:", userData);

      window.dispatchEvent(new Event("storage-update"));

      router.push("/");
      
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <main className="flex-1 flex items-center justify-center px-4 py-16">
        {/* Changed from max-w-md to max-w-lg (512px) */}
        <div className="w-full max-w-lg"> 
          {/* Increased padding from p-8 to p-10 or p-12 */}
          <div className="bg-card rounded-3xl border border-border p-10 shadow-2xl">
            <div className="text-center mb-10">
              <h2 className="text-3xl font-bold text-foreground tracking-tight">
                Log into your account
              </h2>
              <p className="text-base text-muted-foreground mt-3">
                Welcome Back!
              </p>
            </div>
            
            {/* Increased space-y-5 to space-y-6 for more breathing room */}
            {error && (
              <div className="mb-6 flex items-center gap-3 p-4 bg-destructive/10 border border-destructive/20 text-destructive text-sm rounded-2xl animate-in fade-in zoom-in duration-200">
                <AlertCircle className="h-5 w-5" />
                <p>{error}</p>
              </div>
            )}
            
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-3">
                <Label htmlFor="name" className="text-base">Username</Label>
                <Input
                  id="name"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="h-14 rounded-2xl text-lg px-4" // Taller input (h-14) and bigger text
                  required
                />
              </div>

              <div className="space-y-3">
                <Label htmlFor="password" className="text-base">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="h-14 rounded-2xl text-lg px-4 pr-12"
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                  </button>
                </div>
              </div>

              {/* Taller button with larger text */}
              <Button type="submit" className="w-full h-14 rounded-2xl text-lg font-bold shadow-md transition-all hover:scale-[1.01]">
                Login
              </Button>
            </form>

            <p className="text-center text-sm text-muted-foreground mt-8">
              Don't have an account?{" "}
              <Link href="/signup" className="text-primary font-semibold hover:underline">
                Sign Up
              </Link>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}