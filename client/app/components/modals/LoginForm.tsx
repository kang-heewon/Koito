import { login } from "api/api"
import { useState } from "react"
import { AsyncButton } from "../AsyncButton"

export default function LoginForm() {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [remember, setRemember] = useState(false)

    const loginHandler = async () => {
        if (loading) {
            return
        }

        setError('')

        if (!username || !password) {
            if (username || password) {
                setError("username and password are required")
            }
            return
        }

        setLoading(true)

        try {
            const r = await login(username, password, remember)
            if (r.ok) {
                window.location.reload()
                return
            }

            const body = (await r.json().catch(() => null)) as { error?: string } | null
            setError(body?.error ?? `failed to log in (${r.status})`)
        } catch (err) {
            setError(err instanceof Error ? err.message : "failed to log in")
        } finally {
            setLoading(false)
        }
    }

    return (
        <>
        <h2>Log In</h2>
        <div className="flex flex-col items-center gap-4 w-full">
            <p>Logging in gives you access to <strong>admin tools</strong>, such as updating images, merging items, deleting items, and more.</p>
                <form action="#" className="flex flex-col items-center gap-4 w-3/4" onSubmit={(e) => e.preventDefault()}>
                <input
                    name="koito-username"
                    type="text"
                    placeholder="Username"
                    className="w-full mx-auto fg bg rounded p-2"
                    onChange={(e) => setUsername(e.target.value)}
                />
                <input
                    name="koito-password"
                    type="password"
                    placeholder="Password"
                    className="w-full mx-auto fg bg rounded p-2"
                    onChange={(e) => setPassword(e.target.value)}
                />
                <div className="flex gap-2">
                    <input type="checkbox" name="koito-remember" id="koito-remember" onChange={() => setRemember(!remember)} />
                    <label htmlFor="koito-remember">Remember me</label>
                </div>
                <AsyncButton loading={loading} onClick={() => { void loginHandler() }}>Login</AsyncButton>
            </form>
            <p className="error">{error}</p>
        </div>
        </>
    )
}
