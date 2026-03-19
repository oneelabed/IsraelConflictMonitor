"use client"

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff, AlertCircle } from "lucide-react"; // Added AlertCircle icon
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { createUser } from "@/lib/actions/createUser";
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
    setError(""); // Reset error state on new attempt

    // 1. Client-side Validation
    if (username.length < 3) {
      setError("Username must be at least 3 characters long.");
      return;
    }

    if (password.length < 8) {
      setError("Password must be at least 8 characters long.");
      return;
    }
    
    try {
      const userData = await createUser(username.toLowerCase(), password);
      
      // 'expires: 1' means it lasts for 1 day
      Cookies.set("api_key", userData.api_key, { expires: 1, path: '/' });
      Cookies.set("role", userData.role, { expires: 1, path: '/' });
      Cookies.set("username", userData.name, { expires: 1, path: '/' });

      window.dispatchEvent(new Event("storage-update"));
      router.push("/");
      
    } catch (err: any) {
      // 2. Handling Server Errors (e.g., "Username already exists")
      setError(err.message || "Something went wrong during sign up.");
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <main className="flex-1 flex items-center justify-center px-4 py-16">
        <div className="w-full max-w-lg"> 
          <div className="bg-card rounded-3xl border border-border p-10 shadow-2xl">
            <div className="text-center mb-10">
              <h2 className="text-3xl font-bold text-foreground tracking-tight">
                Create your account
              </h2>
              <p className="text-base text-muted-foreground mt-3">
                Join the monitor to get personalized updates
              </p>
            </div>

            {/* ERROR DISPLAY: Shows only if error exists */}
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
                  className={`h-14 rounded-2xl text-lg px-4 ${error.includes("Username") ? "border-destructive" : ""}`}
                  required
                />
              </div>

              <div className="space-y-3">
                <Label htmlFor="password" className="text-base">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="At least 8 characters"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className={`h-14 rounded-2xl text-lg px-4 pr-12 ${error.includes("Password") ? "border-destructive" : ""}`}
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

              <Button type="submit" className="w-full h-14 rounded-2xl text-lg font-bold shadow-md transition-all hover:scale-[1.01]">
                Sign Up
              </Button>
            </form>

            <p className="text-center text-sm text-muted-foreground mt-8">
              Already have an account?{" "}
              <Link href="/login" className="text-primary font-semibold hover:underline">
                Log in
              </Link>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}