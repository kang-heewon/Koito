import { logout, updateUser } from "api/api"
import { useState } from "react"
import { AsyncButton } from "../AsyncButton"
import { useAppContext } from "~/providers/AppProvider"

export default function Account() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [confirmPw, setConfirmPw] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const [success, setSuccess] = useState('')
    const { user, setUsername: setCtxUsername } = useAppContext()

    const logoutHandler = async () => {
        if (loading) {
            return
        }

        setError('')
        setSuccess('')
        setLoading(true)

        try {
            const r = await logout()
            if (r.ok) {
                window.location.reload()
            } else {
                const body = (await r.json().catch(() => null)) as { error?: string } | null
                setError(body?.error ?? 'failed to log out')
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'failed to log out')
        } finally {
            setLoading(false)
        }
    }

    const updateHandler = async () => {
        if (loading) {
            return
        }

        setError('')
        setSuccess('')

        if (username.trim() === '' && password === '') {
            setError('provide a new username or password before submitting')
            return
        }

        if (password !== "" && confirmPw === "") {
            setError("confirm your new password before submitting")
            return
        }

        if (password !== confirmPw && confirmPw !== '') {
            setError('new password and confirmation must match')
            return
        }

        setLoading(true)

        try {
            const r = await updateUser(username, password)
            if (r.ok) {
                setSuccess("sucessfully updated user")
                if (username !== "") {
                    setCtxUsername(username)
                }
                setUsername('')
                setPassword('')
                setConfirmPw('')
            } else {
                const body = (await r.json().catch(() => null)) as { error?: string } | null
                setError(body?.error ?? 'failed to update user')
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'failed to update user')
        } finally {
            setLoading(false)
        }
    }

    return (
        <>
        <h2>Account</h2>
        <div className="flex flex-col gap-6">
            <div className="flex flex-col gap-4 items-center">
                <p>You're logged in as <strong>{user?.username}</strong></p>
                <AsyncButton loading={loading} onClick={logoutHandler}>Logout</AsyncButton>
            </div>
            <h2>Update User</h2>
            <form action="#" onSubmit={(e) => e.preventDefault()} className="flex flex-col gap-4">
                <div className="flex flex gap-4">
                    <input
                        name="koito-update-username"
                        type="text"
                        placeholder="Update username"
                        className="w-full mx-auto fg bg rounded p-2"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>
                <div className="w-sm">
                    <AsyncButton loading={loading} onClick={updateHandler}>Submit</AsyncButton>
                </div>
            </form>
            <form action="#" onSubmit={(e) => e.preventDefault()} className="flex flex-col gap-4">
                <div className="flex flex gap-4">
                    <input
                        name="koito-update-password"
                        type="password"
                        placeholder="Update password"
                        className="w-full mx-auto fg bg rounded p-2"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <input
                        name="koito-confirm-password"
                        type="password"
                        placeholder="Confirm new password"
                        className="w-full mx-auto fg bg rounded p-2"
                        value={confirmPw}
                        onChange={(e) => setConfirmPw(e.target.value)}
                    />
                </div>
                <div className="w-sm">
                    <AsyncButton loading={loading} onClick={updateHandler}>Submit</AsyncButton>
                </div>
            </form>
            {success !== "" && <p className="success">{success}</p>}
            {error !== "" && <p className="error">{error}</p>}
        </div>
        </>
    )
}
