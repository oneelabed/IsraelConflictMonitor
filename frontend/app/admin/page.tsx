"use client"

import { useEffect, useState } from "react"
import { Plus, Globe, Link as LinkIcon, Image as ImageIcon, Loader2, Users, ShieldCheck, Calendar, Key } from "lucide-react"
import { useRouter } from "next/navigation";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function AdminPage() {
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState({ text: "", type: "" })
  const [formData, setFormData] = useState({ name: "", url: "", icon_url: "" })
  const [users, setUsers] = useState<any[]>([])
  const [loadingUsers, setLoadingUsers] = useState(true)
  const router = useRouter();
  
  useEffect(() => {
    const apiKey = localStorage.getItem("api_key");
    const role = localStorage.getItem("role");

    if (!apiKey || role !== "admin") {
      // If not an admin, send them to the home page or a "Forbidden" page
      router.push("/");
    }
  }, [router]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await fetch(`${API_URL}/v1/admin/users`, {
          headers: { 'Authorization': `ApiKey ${localStorage.getItem("api_key")}` }
        })
        const data = await res.json()
        setUsers(data)
      } catch (err) {
        console.error("Failed to load users", err)
      } finally {
        setLoadingUsers(false)
      }
    }
    fetchUsers()
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    try {
      const res = await fetch(`${API_URL}/v1/feeds`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `ApiKey ${localStorage.getItem("api_key")}` 
        },
        body: JSON.stringify(formData)
      })

      if (!res.ok) throw new Error("Unauthorized or invalid data")

      setMessage({ text: "Feed added successfully!", type: "success" })
      setFormData({ name: "", url: "", icon_url: "" })
    } catch (err) {
      setMessage({ text: "Failed to add feed. Are you an admin?", type: "error" })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Admin Dashboard</h1>
      
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-2xl border shadow-sm space-y-4">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <Plus className="w-5 h-5" /> Add New RSS Source
        </h2>

        {message.text && (
          <div className={`p-3 rounded-md text-sm ${message.type === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
            {message.text}
          </div>
        )}

        <div className="space-y-2">
          <label className="text-sm font-medium">Source Name (e.g., Walla News)</label>
          <div className="relative">
            <Globe className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
            <input 
              required
              className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary/20 outline-none"
              placeholder="Source Name"
              value={formData.name}
              onChange={e => setFormData({...formData, name: e.target.value})}
            />
          </div>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">RSS Feed URL</label>
          <div className="relative">
            <LinkIcon className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
            <input 
              required
              type="url"
              className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary/20 outline-none"
              placeholder="https://example.com/rss"
              value={formData.url}
              onChange={e => setFormData({...formData, url: e.target.value})}
            />
          </div>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">Icon URL</label>
          <div className="relative">
            <ImageIcon className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
            <input 
              required
              className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary/20 outline-none"
              placeholder="https://example.com/logo.png"
              value={formData.icon_url}
              onChange={e => setFormData({...formData, icon_url: e.target.value})}
            />
          </div>
        </div>

        <button 
          disabled={loading}
          type="submit"
          className="w-full bg-slate-900 text-white py-3 rounded-lg font-bold hover:bg-slate-800 transition-all disabled:opacity-50"
        >
          {loading ? <Loader2 className="animate-spin mx-auto" /> : "Register New Source"}
        </button>
      </form>
      <div className="mt-12">
        <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Users className="w-6 h-6 text-primary" /> Registered Users
        </h2>

        <div className="bg-white rounded-2xl border shadow-sm overflow-hidden">
            {loadingUsers ? (
                <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                    <Loader2 className="animate-spin h-8 w-8 mb-2" />
                    <p className="text-sm">Fetching user database...</p>
                </div>
            ) : (
                <table className="w-full text-left border-collapse">
                    <thead className="bg-slate-50 border-b">
                        <tr>
                            <th className="px-6 py-4 text-sm font-semibold text-slate-600">User</th>
                            <th className="px-6 py-4 text-sm font-semibold text-slate-600">Joined</th>
                            <th className="px-6 py-4 text-sm font-semibold text-slate-600">API Key</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y">
                        {users.map((u) => (
                        <tr key={u.id} className="hover:bg-slate-50 transition-colors">
                            <td className="px-6 py-4">
                            <span className="font-bold text-slate-900">{u.name}</span>
                            <div className="text-xs text-slate-400 font-mono">{u.id}</div>
                            </td>
                            <td className="px-6 py-4 text-sm text-slate-600">
                            <div className="flex items-center gap-1">
                                <Calendar className="w-3.5 h-3.5" />
                                {new Date(u.created_at).toLocaleDateString()}
                            </div>
                            </td>
                            <td className="px-6 py-4">
                            <code className="text-[10px] bg-slate-100 p-1 rounded select-all cursor-pointer">
                                {u.api_key.substring(0, 8)}...
                            </code>
                            </td>
                        </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
      </div>
    </div>
  )
}