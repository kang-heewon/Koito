import { useState } from "react";
import { Modal } from "./Modal";
import { AsyncButton } from "../AsyncButton";
import { useSubmitListen } from "../../hooks/useSubmitListen";
import { useNavigate } from "react-router";

interface Props {
    open: boolean 
    setOpen: Function
    trackid: number
}

export default function AddListenModal({ open, setOpen, trackid }: Props) {
    const [ts, setTS] = useState<Date>(new Date);
    const [error, setError] = useState('')
    const navigate = useNavigate()
    const submitMutation = useSubmitListen()
    const close = () => {
        setOpen(false)
    }
    const submit = async () => {
        setError('')
        try {
            const r = await submitMutation.mutateAsync({ trackId: trackid.toString(), ts })
            if (r.ok) {
                navigate(0)
                return
            }
            const body = (await r.json().catch(() => null)) as { error?: string } | null
            setError(body?.error ?? `failed to add listen (${r.status})`)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'failed to add listen')
        }
    }

    const formatForDatetimeLocal = (d: Date) => {
        const pad = (n: number) => n.toString().padStart(2, "0");
        return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
    };

    return (
        <Modal isOpen={open} onClose={close}>
            <h2>Add Listen</h2>
            <div className="flex flex-col items-center gap-4">
                <input
                    type="datetime-local"
                    className="w-full mx-auto fg bg rounded p-2"
                    value={formatForDatetimeLocal(ts)}
                    onChange={(e) => setTS(new Date(e.target.value))}
                />
                <AsyncButton loading={submitMutation.isPending} onClick={() => { void submit() }}>Submit</AsyncButton>
                <p className="error">{error}</p>
            </div>
        </Modal>
    )
}
